package route

import "opennamu/route/tool"

func Api_add_user_hash(config tool.Config, id string, password_hash string, email string, encode string) map[string]any {
	return Api_add_user_hash_internal(config, id, password_hash, email, encode, "", false)
}

func Api_add_user_hash_invite(config tool.Config, id string, password_hash string, email string, encode string, invite_hash string) map[string]any {
	return Api_add_user_hash_internal(config, id, password_hash, email, encode, invite_hash, true)
}
