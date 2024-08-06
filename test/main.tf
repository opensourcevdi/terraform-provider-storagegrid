resource "storagegrid_user" "a" {
  name = "user1"
  # Computed:
  # id
  # user_urn
}

resource "storagegrid_access_key" "key1" {
  user_id = storagegrid_user.a.id
  # Computed:
  # id
  # user_urn
  # access_key_id
  # secret_access_key
}

resource "storagegrid_access_key" "key2" {
  user_id = storagegrid_user.a.id
}
