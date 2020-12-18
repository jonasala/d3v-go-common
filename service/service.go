package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	consul "github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
)

//Service agrega as configurações de banco de dados (pgsql) e service discovery (consul)
type Service struct {
	ID             string
	Name           string
	Config         Config
	ConsulAgent    *consul.Agent
	ConsulKV       *consul.KV
	RedisClient    *redis.Client
	SpecificConfig map[string]interface{}
}

//Config representa as configurações de que vem do consul
type Config struct {
	DBHost             string `json:"db_host"`
	DBName             string `json:"db_name"`
	DBUser             string `json:"db_user"`
	DBPassword         string `json:"db_password"`
	DBPort             string `json:"db_port"`
	DBSchema           string `json:"db_schema"`
	RedisServer        string `json:"redis_server"`
	HTTPPort           string `json:"http_port"`
	HTTPAddress        string `json:"http_address"`
	FabioAddress       string `json:"fabio_address"`
	TTL                string `json:"ttl"`
	JWTSecret          string `json:"jwt_secret"`
	RefreshTokenSecret string `json:"refresh_token_secret"`
}

//New cria e configura um serviço
func New(name, configKey string) (*Service, error) {
	service := &Service{
		ID:   uuid.New(),
		Name: name,
	}
	consulConfig := consul.DefaultConfig()
	if consulAddr := os.Getenv("CONSUL_ADDR"); consulAddr != "" {
		consulConfig.Address = consulAddr
	}

	client, err := consul.NewClient(consulConfig)
	if err != nil {
		return service, fmt.Errorf("não foi possível conectar com o consul. %v", err)
	}

	service.ConsulKV = client.KV()

	pair, _, err := service.ConsulKV.Get(configKey, nil)
	if err != nil {
		return service, fmt.Errorf("problemas ao recuperar a chave %v. %v", configKey, err)
	}

	if pair == nil {
		return service, fmt.Errorf("chave %v inexistente", configKey)
	}

	config := Config{}
	if err := json.Unmarshal(pair.Value, &config); err != nil {
		return service, fmt.Errorf("a chave %v não é um json válido. %v", configKey, err)
	}

	service.Config = config

	//configurações específicas do serviço
	pair, _, err = service.ConsulKV.Get("config/"+name, nil)
	if err != nil {
		log.Printf("problemas ao recuperar a chave de config específica %v. %v\n", name, err)
	} else if pair == nil {
		log.Printf("chave de config específica %v inexistente", name)
	} else if err := json.Unmarshal(pair.Value, &service.SpecificConfig); err != nil {
		log.Printf("a chave de config específica %v não é um json válido. %v", name, err)
	}

	if mode := os.Getenv("DOCKER_MODE"); mode == "dev" {
		service.Config.HTTPAddress = "host.docker.internal"
	} else {
		ip, err := IPAddr()
		if err != nil {
			return service, fmt.Errorf("não foi possível determinar o ip para registrar este serviço. %v", err)
		}
		service.Config.HTTPAddress = ip.String()
	}

	if os.Getenv("HTTP_PORT") != "" {
		service.Config.HTTPPort = os.Getenv("HTTP_PORT")
	}
	if service.Config.HTTPPort == "" {
		service.Config.HTTPPort = "80"
	}

	service.ConsulAgent = client.Agent()

	if service.Config.RedisServer != "" {
		service.RedisClient = redis.NewClient(&redis.Options{
			Addr:     service.Config.RedisServer,
			Password: "",
			DB:       0,
		})

		if _, err := service.RedisClient.Ping().Result(); err != nil {
			return service, fmt.Errorf("não foi possível estabelecer conexão com o redis. %v", err)
		}
	}

	return service, nil
}

//RegisterService registra o serviço no consul.
//O parâmetro healthCheckFunction deve ser uma função que retorna um erro se houver falha
func (s *Service) RegisterService(healthcheckFunction func() error) error {
	ttl, err := time.ParseDuration(s.Config.TTL)
	if err != nil {
		return fmt.Errorf("não foi possível interpretar o ttl. %v", err)
	}

	port, _ := strconv.Atoi(s.Config.HTTPPort)

	def := &consul.AgentServiceRegistration{
		ID:      s.ID,
		Name:    s.Name,
		Address: s.Config.HTTPAddress,
		Port:    port,
		Tags:    []string{"urlprefix-/" + s.Name + " strip=/" + s.Name},
		Check: &consul.AgentServiceCheck{
			TTL: ttl.String(),
		},
	}

	if err = s.ConsulAgent.ServiceRegister(def); err != nil {
		return fmt.Errorf("não foi possível registrar o serviço no consul. %v", err)
	}

	go func() {
		ticker := time.NewTicker(ttl / 2)
		for range ticker.C {
			status := "pass"
			output := ""
			err := healthcheckFunction()
			if err != nil {
				status = "fail"
				output = err.Error()
			}
			s.ConsulAgent.UpdateTTL("service:"+s.ID, output, status)
		}
	}()

	return nil
}

//GracefullyShutdown faz o desregistramento no consul (não funciona com watcher, não registre)
func (s *Service) GracefullyShutdown() {
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sign
		s.ConsulAgent.ServiceDeregister(s.ID)
		os.Exit(1)
	}()
}

//IPAddr recupera o endereco IP do serviço
func IPAddr() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			if ipnet.IP.To4() != nil || ipnet.IP.To16() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, nil
}
