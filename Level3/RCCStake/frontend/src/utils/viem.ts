import { holesky, sepolia } from "viem/chains";
import { PublicClient, createPublicClient, http } from 'viem'
import { mainnet } from 'viem/chains'





export const viemClients = (chaiId: number): PublicClient => {
  const clients: {
    [key: number]: PublicClient
  } = {
    [holesky.id]: createPublicClient({
      chain: holesky,
      transport: http('https://api.zan.top/node/v1/eth/holesky/XXXXXXXXXXXXXXXX')
    })
  }
  return clients[chaiId]
}