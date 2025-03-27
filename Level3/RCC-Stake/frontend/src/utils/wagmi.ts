import { getDefaultConfig } from '@rainbow-me/rainbowkit';
import { http } from 'viem';
import {
  arbitrum,
  base,
  holesky,
  mainnet,
  optimism,
  polygon,
  sepolia,
} from 'wagmi/chains';
// from https://cloud.walletconnect.com/
const ProjectId = 'e3242412afd6123ce1dda1de23a8c016'

export const config = getDefaultConfig({
  appName: 'Rcc Stake',
  projectId: ProjectId,
  chains: [
    holesky
  ],
  transports: {
    // 替换之前 不可用的 https://rpc.sepolia.org/
    [holesky.id]: http('https://api.zan.top/node/v1/eth/holesky/XXXXXXXXXXX')
  },
  ssr: true,
});

export const defaultChainId: number = holesky.id