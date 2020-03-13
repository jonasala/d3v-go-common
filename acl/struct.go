package acl

//ACL representa as permissões a nível de registro x usuário
type ACL struct {
	CanView   bool `json:"can_view"`
	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
}

//NewACLAllowAll cria uma permissão ACL permitindo todas as ações
func NewACLAllowAll() *ACL {
	return &ACL{
		CanView:   true,
		CanEdit:   true,
		CanDelete: true,
	}
}
