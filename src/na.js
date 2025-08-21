//
// NetAngels API
//
const AUTH = "https://panel.netangels.ru/api/gateway/token/"
const API = 'https://api-ms.netangels.ru/api/v1/dns/'

const auth = doAuth()

export async function doAuth(key = process.env.NETANGELS_API_KEY) {
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
    headers: {
      authorization: `Bearer ${j.token}`
    }
  }
}

export async function zones() {
  let q = await fetch(`${API}zones`, await auth)
  return await q.json()
}
