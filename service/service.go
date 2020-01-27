package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"
)

//Service agrega as configurações de banco de dados (pgsql) e service discovery (consul)
type Service struct {
	Name        string
	Config      Config
	ConsulAgent *consul.Agent
	ConsulKV    *consul.KV
}

//Config representa as configurações de que vem do consul
type Config struct {
	DBHost      string `json:"db_host"`
	DBName      string `json:"db_name"`
	DBUser      string `json:"db_user"`
	DBPassword  string `json:"db_password"`
	DBPort      string `json:"db_port"`
	DBSchema    string `json:"db_schema"`
	HTTPPort    string `json:"http_port"`
	HTTPAddress string `json:"http_address"`
	TTL         string `json:"ttl"`
}

//New cria e configura um serviço
func New(name, configKey string) (*Service, error) {
	service := &Service{Name: name}

	client, err := consul.NewClient(consul.DefaultConfig())
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

	service.ConsulAgent = client.Agent()

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
		Name:    s.Name,
		Address: s.Config.HTTPAddress,
		Port:    port,
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
			err := healthcheckFunction()
			if err != nil {
				s.ConsulAgent.FailTTL("service:"+s.Name, err.Error())
			} else {
				s.ConsulAgent.PassTTL("service:"+s.Name, "")
			}
		}
	}()

	return nil
}
