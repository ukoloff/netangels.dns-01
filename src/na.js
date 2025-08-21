//
// NetAngels API
//
const AUTH = "https://panel.netangels.ru/api/gateway/token/"

export async function auth(key = process.env.NETANGELS_API_KEY) {
  let params = new FormData()
  params.append('api_key', key)
  let auth = await fetch(AUTH, {
    method: 'POST',
    body: params,
  })
  let j = await auth.json()
  if (!j.token) {
    throw new Error('Failed to authorize to netangels.ru')
  }
  return {
    authorization: j.token
  }
}
