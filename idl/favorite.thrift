namespace go favorite
include "model.thrift"

struct GetFavoriteReq {
  required string type
}
struct GetFavoriteResp {
  required  model.BaseResp resp
  required list<model.Favorite> items
}

struct AddFavoriteReq {
  required i64 target_id
  required string target_type
}
struct AddFavoriteResp {
  required  model.BaseResp resp
  required  model.Favorite favorite
}

struct RemoveFavoriteReq {
  required i64 favorite_id
}
struct RemoveFavoriteResp {
  required  model.BaseResp resp
}

service FavoriteService {
    GetFavoriteResp GetFavorites(1: GetFavoriteReq req)(api.get="/api/users/me/favorites")
    AddFavoriteResp AddFavorite(1: AddFavoriteReq req)(api.post="/api/users/me/favorites")
    RemoveFavoriteResp RemoveFavorite(1: RemoveFavoriteReq req)(api.delete="/api/users/me/favorites")
}
