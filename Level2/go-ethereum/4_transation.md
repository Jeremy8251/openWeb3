## transaction

### 源码

https://github.com/ethereum/go-ethereum/blob/release/1.15/core/types/transaction.go

### 交易流程

1. 发起交易

2. 交易签名
3. 提交交易
4. 广播交易：通知EVM执行，同时把交易信息广播给其他节点

https://github.com/ethereum/go-ethereum/blob/release/1.15/internal/ethapi/api.go

https://github.com/ethereum/go-ethereum/blob/release/1.15/internal/ethapi/transaction_args.go

调用链

```text
JSON-RPC 请求 → eth/api.go 接收 → TransactionAPI.SendTransaction() → 
    |-- 账户查找 → accounts/manager.go:Find()
    |-- Nonce 锁定 → eth/api.go:nonceLock
    |-- 参数默认值 → core/transaction_args.go:setDefaults()
    |-- 交易创建 → core/types/transaction.go:NewTx()
    |-- 签名 → accounts/wallet.go:SignTx()
    |-- 提交交易 → eth/api.go:SubmitTransaction() → 
        |-- 交易池验证 → core/txpool/pool.go:add()
        |-- 广播 → eth/handler.go:BroadcastTx()

```

internal/ethapi/api.go

创建一个新的交易

```go
// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
// 通过给定的参数创建一个新的交易并且提交到txpool
// 用户通过JSON RPC发起的交易web.eth.SendTransaction(...)
func (api *TransactionAPI) SendTransaction(ctx context.Context, args TransactionArgs) (common.Hash, error) {
	// Look up the wallet containing the requested signer
    // 查找一个包含了所需要的签名的钱包地址
	account := accounts.Account{Address: args.from()}
	// 通过账号查找对应的钱包结构
	wallet, err := api.b.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}

	if args.Nonce == nil {
		// Hold the mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		api.nonceLock.LockAddr(args.from())
		defer api.nonceLock.UnlockAddr(args.from())
	}
	if args.IsEIP4844() {
		return common.Hash{}, errBlobTxNotSupported
	}

	// Set some sanity defaults and terminate on failure
    // 设置默认值
	if err := args.setDefaults(ctx, api.b, false); err != nil {
		return common.Hash{}, err
	}
	// Assemble the transaction and sign with the wallet
    // 1.通过wallet创建transaction，最终到transaction创建交易
 // [transaction_args.go] ToTransaction-》[transaction.go] types.NewTx(data)
	tx := args.ToTransaction(types.LegacyTxType)
    // 2.对交易进行签名
	signed, err := wallet.SignTx(account, tx, api.b.ChainConfig().ChainID)
	if err != nil {
		return common.Hash{}, err
	}
    //3.提交交易 
	return SubmitTransaction(ctx, api.b, signed)
}
```

transaction_args.go：发起交易所需要的参数

```go
// TransactionArgs represents the arguments to construct a new transaction
// or a message call.
// 发起交易所需要的参数
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`

	// For BlobTxType
	BlobFeeCap *hexutil.Big  `json:"maxFeePerBlobGas"`
	BlobHashes []common.Hash `json:"blobVersionedHashes,omitempty"`

	// For BlobTxType transactions with blob sidecar
	Blobs       []kzg4844.Blob       `json:"blobs"`
	Commitments []kzg4844.Commitment `json:"commitments"`
	Proofs      []kzg4844.Proof      `json:"proofs"`

	// For SetCodeTxType
	AuthorizationList []types.SetCodeAuthorization `json:"authorizationList"`

	// This configures whether blobs are allowed to be passed.
	blobSidecarAllowed bool
}

// ToTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
// 设置参数，最后向transaction发起一个交易
func (args *TransactionArgs) ToTransaction(defaultType int) *types.Transaction {
	usedType := types.LegacyTxType
	switch {
	case args.AuthorizationList != nil || defaultType == types.SetCodeTxType:
		usedType = types.SetCodeTxType
	case args.BlobHashes != nil || defaultType == types.BlobTxType:
		usedType = types.BlobTxType
	case args.MaxFeePerGas != nil || defaultType == types.DynamicFeeTxType:
		usedType = types.DynamicFeeTxType
	case args.AccessList != nil || defaultType == types.AccessListTxType:
		usedType = types.AccessListTxType
	}
	// Make it possible to default to newer tx, but use legacy if gasprice is provided
	if args.GasPrice != nil {
		usedType = types.LegacyTxType
	}
	///....................
    // 发起交易
	return types.NewTx(data)
}
```

软件钱包accounts/keystore/ wallet.go和keystore.go

```go
// wallet.go
func (w *keystoreWallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTx(account, tx, chainID)
}

//keystore.go
// SignTx signs the given transaction with the requested account.
func (ks *KeyStore) SignTx(a accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Look up the key to sign with and abort if it cannot be found
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	unlockedKey, found := ks.unlocked[a.Address]
	if !found {
		return nil, ErrLocked
	}
	// Depending on the presence of the chain ID, sign with 2718 or homestead
	signer := types.LatestSignerForChainID(chainID)
	return types.SignTx(tx, signer, unlockedKey.PrivateKey)
}
```

硬件钱包accounts/usbwallet/wallet.go

调用链

```text
用户发起交易 → RPC调用SendTransaction → eth/api_backend.go → 账号管理器选择硬件钱包 → accounts/usbwallet/wallet.go → SignTx → 通过USB HID与设备通信 → 返回签名结果
```

```go
// 使用硬件钱包（如Ledger）对交易进行签名，需用户物理确认,参数：要签名的账户,待签名的原始交易对象区块链网络ID
func (w *wallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
    // 防止在签名过程中钱包状态发生变化
	w.stateLock.RLock() // Comms have own mutex, this is for the state fields
	defer w.stateLock.RUnlock()

	// If the wallet is closed, abort
    // 设备连接检查
	if w.device == nil {
		return nil, accounts.ErrWalletClosed
	}
	// Make sure the requested account is contained within
    // 账户路径查询
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}
	// All infos gathered and metadata checks out, request signing
    // 使用缓冲大小为1的channel实现互斥锁
	<-w.commsLock
	defer func() { w.commsLock <- struct{}{} }()

	// Ensure the device isn't screwed with while user confirmation is pending
	// TODO(karalabe): remove if hotplug lands on Windows
	w.hub.commsLock.Lock()
	w.hub.commsPend++
	w.hub.commsLock.Unlock()

	defer func() {
		w.hub.commsLock.Lock()
		w.hub.commsPend--
		w.hub.commsLock.Unlock()
	}()
	// Sign the transaction and verify the sender to avoid hardware fault surprises
    // 通过驱动（如ledgerDriver）发送APDU指令到设备
	sender, signed, err := w.driver.SignTx(path, tx, chainID)
	if err != nil {
		return nil, err
	}
	if sender != account.Address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", account.Address.Hex(), sender.Hex())
	}
	return signed, nil
}
```

**api.go/SubmitTransaction 入口**‌

文件：`eth/api.go`

```go
func SubmitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
    // 调用链核心入口
    if err := b.SendTx(ctx, tx); err != nil {
        return common.Hash{}, err
    }
    return tx.Hash(), nil
}

```

**Backend.SendTx() 实现**‌

文件：`eth/backend.go`

```go
func (b *Ethereum) SendTx(ctx context.Context, tx *types.Transaction) error {
    // 提交到交易池
    return b.txPool.Add([]*types.Transaction{tx}, false, false)
}

```

**TxPool.Add() 入口**

```go
func (p *TxPool) Add(txs []*types.Transaction, local, sync bool) []error {
    // 批量添加交易
    return p.add(txs, local, sync)
}

```

**TxPool.add() 核心逻辑**‌

文件：`core/txpool/pool.go`

```go
func (p *TxPool) add(txs []*types.Transaction, local, sync bool) []error {
    // 交易验证（Nonce、签名、Gas等）
    for _, tx := range txs {
        if err := p.validateTx(tx, local); err != nil {
            errors = append(errors, err)
            continue
        }
        // 加入待处理队列
        replace, err := p.enqueueTx(tx.Hash(), tx)
    }

    // 触发新交易事件（关键广播入口）
    if len(newTxs) > 0 {
        go p.txFeed.Send(core.NewTxsEvent{newTxs})
    }
}

```

 **P2P 广播触发**‌

文件：`eth/handler.go`

````go
// 在 handler 初始化时订阅交易池事件
func NewHandler(...) *Handler {
    handler.txFetcher.Start()
    handler.txSub = handler.txPool.SubscribeNewTxsEvent(handler.txFetcher.Notify)
}

// 事件监听回调
func (f *TxFetcher) Notify(event core.NewTxsEvent) {
    for _, tx := range event.Txs {
        // 标记需要广播的交易
        f.announced[tx.Hash()] = struct{}{}
        // 触发广播
        f.broadcastTx(tx)
    }
}

// 实际广播逻辑
func (h *handler) broadcastTx(tx *types.Transaction) {
    // 封装为 P2P 消息
    msg := eth.NewMessage(eth.TxMsg, []*types.Transaction{tx})
    // 发送给所有连接的 peer
    h.peers.broadcast(msg)
}

````

### transaction.go源码

func NewTx(inner TxData) *Transaction

```go
package types

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrInvalidSig           = errors.New("invalid transaction v, r, s values")
	ErrUnexpectedProtection = errors.New("transaction type does not supported EIP-155 protected signatures")
	ErrInvalidTxType        = errors.New("transaction type not valid in this context")
	ErrTxTypeNotSupported   = errors.New("transaction type not supported")
	ErrGasFeeCapTooLow      = errors.New("fee cap less than base fee")
	errShortTypedTx         = errors.New("typed transaction too short")
	errInvalidYParity       = errors.New("'yParity' field must be 0 or 1")
	errVYParityMismatch     = errors.New("'v' and 'yParity' fields do not match")
	errVYParityMissing      = errors.New("missing 'yParity' or 'v' field in transaction")
)

// Transaction types.
const (
	LegacyTxType     = 0x00
	AccessListTxType = 0x01
	DynamicFeeTxType = 0x02
	BlobTxType       = 0x03
	SetCodeTxType    = 0x04
)

// Transaction is an Ethereum transaction.
type Transaction struct {
	inner TxData    // Consensus contents of a transaction
	time  time.Time // Time first seen locally (spam avoidance)

	// caches
	hash atomic.Pointer[common.Hash]
	size atomic.Uint64
	from atomic.Pointer[sigCache]
}

// NewTx creates a new transaction.
// 创建新的交易
func NewTx(inner TxData) *Transaction {
	tx := new(Transaction)
	tx.setDecoded(inner.copy(), 0)
	return tx
}

// TxData is the underlying data of a transaction.
//
// This is implemented by DynamicFeeTx, LegacyTx and AccessListTx.
type TxData interface {
	txType() byte // returns the type ID
	copy() TxData // creates a deep copy and initializes all fields

	chainID() *big.Int
	accessList() AccessList
	data() []byte
	gas() uint64
	gasPrice() *big.Int
	gasTipCap() *big.Int
	gasFeeCap() *big.Int
	value() *big.Int
	nonce() uint64
	to() *common.Address

	rawSignatureValues() (v, r, s *big.Int)
	setSignatureValues(chainID, v, r, s *big.Int)

	// effectiveGasPrice computes the gas price paid by the transaction, given
	// the inclusion block baseFee.
	//
	// Unlike other TxData methods, the returned *big.Int should be an independent
	// copy of the computed value, i.e. callers are allowed to mutate the result.
	// Method implementations can use 'dst' to store the result.
	effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int

	encode(*bytes.Buffer) error
	decode([]byte) error
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	if tx.Type() == LegacyTxType {
		return rlp.Encode(w, tx.inner)
	}
	// It's an EIP-2718 typed TX envelope.
	buf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(buf)
	buf.Reset()
	if err := tx.encodeTyped(buf); err != nil {
		return err
	}
	return rlp.Encode(w, buf.Bytes())
}

// encodeTyped writes the canonical encoding of a typed transaction to w.
func (tx *Transaction) encodeTyped(w *bytes.Buffer) error {
	w.WriteByte(tx.Type())
	return tx.inner.encode(w)
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (tx *Transaction) MarshalBinary() ([]byte, error) {
	if tx.Type() == LegacyTxType {
		return rlp.EncodeToBytes(tx.inner)
	}
	var buf bytes.Buffer
	err := tx.encodeTyped(&buf)
	return buf.Bytes(), err
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	kind, size, err := s.Kind()
	switch {
	case err != nil:
		return err
	case kind == rlp.List:
		// It's a legacy transaction.
		var inner LegacyTx
		err := s.Decode(&inner)
		if err == nil {
			tx.setDecoded(&inner, rlp.ListSize(size))
		}
		return err
	case kind == rlp.Byte:
		return errShortTypedTx
	default:
		// It's an EIP-2718 typed TX envelope.
		// First read the tx payload bytes into a temporary buffer.
		b, buf, err := getPooledBuffer(size)
		if err != nil {
			return err
		}
		defer encodeBufferPool.Put(buf)
		if err := s.ReadBytes(b); err != nil {
			return err
		}
		// Now decode the inner transaction.
		inner, err := tx.decodeTyped(b)
		if err == nil {
			tx.setDecoded(inner, size)
		}
		return err
	}
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP-2718 typed transactions.
func (tx *Transaction) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		var data LegacyTx
		err := rlp.DecodeBytes(b, &data)
		if err != nil {
			return err
		}
		tx.setDecoded(&data, uint64(len(b)))
		return nil
	}
	// It's an EIP-2718 typed transaction envelope.
	inner, err := tx.decodeTyped(b)
	if err != nil {
		return err
	}
	tx.setDecoded(inner, uint64(len(b)))
	return nil
}

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *Transaction) decodeTyped(b []byte) (TxData, error) {
	if len(b) <= 1 {
		return nil, errShortTypedTx
	}
	var inner TxData
	switch b[0] {
	case AccessListTxType:
		inner = new(AccessListTx)
	case DynamicFeeTxType:
		inner = new(DynamicFeeTx)
	case BlobTxType:
		inner = new(BlobTx)
	case SetCodeTxType:
		inner = new(SetCodeTx)
	default:
		return nil, ErrTxTypeNotSupported
	}
	err := inner.decode(b[1:])
	return inner, err
}

// setDecoded sets the inner transaction and size after decoding.
func (tx *Transaction) setDecoded(inner TxData, size uint64) {
	tx.inner = inner
	tx.time = time.Now()
	if size > 0 {
		tx.size.Store(size)
	}
}

func sanityCheckSignature(v *big.Int, r *big.Int, s *big.Int, maybeProtected bool) error {
	if isProtectedV(v) && !maybeProtected {
		return ErrUnexpectedProtection
	}

	var plainV byte
	if isProtectedV(v) {
		chainID := deriveChainId(v).Uint64()
		plainV = byte(v.Uint64() - 35 - 2*chainID)
	} else if maybeProtected {
		// Only EIP-155 signatures can be optionally protected. Since
		// we determined this v value is not protected, it must be a
		// raw 27 or 28.
		plainV = byte(v.Uint64() - 27)
	} else {
		// If the signature is not optionally protected, we assume it
		// must already be equal to the recovery id.
		plainV = byte(v.Uint64())
	}
	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return ErrInvalidSig
	}

	return nil
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28 && v != 1 && v != 0
	}
	// anything not 27 or 28 is considered protected
	return true
}

// Protected says whether the transaction is replay-protected.
func (tx *Transaction) Protected() bool {
	switch tx := tx.inner.(type) {
	case *LegacyTx:
		return tx.V != nil && isProtectedV(tx.V)
	default:
		return true
	}
}

// Type returns the transaction type.
func (tx *Transaction) Type() uint8 {
	return tx.inner.txType()
}

// ChainId returns the EIP155 chain ID of the transaction. The return value will always be
// non-nil. For legacy transactions which are not replay-protected, the return value is
// zero.
func (tx *Transaction) ChainId() *big.Int {
	return tx.inner.chainID()
}

// Data returns the input data of the transaction.
func (tx *Transaction) Data() []byte { return tx.inner.data() }

// AccessList returns the access list of the transaction.
func (tx *Transaction) AccessList() AccessList { return tx.inner.accessList() }

// Gas returns the gas limit of the transaction.
func (tx *Transaction) Gas() uint64 { return tx.inner.gas() }

// GasPrice returns the gas price of the transaction.
func (tx *Transaction) GasPrice() *big.Int { return new(big.Int).Set(tx.inner.gasPrice()) }

// GasTipCap returns the gasTipCap per gas of the transaction.
func (tx *Transaction) GasTipCap() *big.Int { return new(big.Int).Set(tx.inner.gasTipCap()) }

// GasFeeCap returns the fee cap per gas of the transaction.
func (tx *Transaction) GasFeeCap() *big.Int { return new(big.Int).Set(tx.inner.gasFeeCap()) }

// Value returns the ether amount of the transaction.
func (tx *Transaction) Value() *big.Int { return new(big.Int).Set(tx.inner.value()) }

// Nonce returns the sender account nonce of the transaction.
func (tx *Transaction) Nonce() uint64 { return tx.inner.nonce() }

// To returns the recipient address of the transaction.
// For contract-creation transactions, To returns nil.
func (tx *Transaction) To() *common.Address {
	return copyAddressPtr(tx.inner.to())
}

// Cost returns (gas * gasPrice) + (blobGas * blobGasPrice) + value.
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas()))
	if tx.Type() == BlobTxType {
		total.Add(total, new(big.Int).Mul(tx.BlobGasFeeCap(), new(big.Int).SetUint64(tx.BlobGas())))
	}
	total.Add(total, tx.Value())
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
// The return values may be nil or zero, if the transaction is unsigned.
func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	return tx.inner.rawSignatureValues()
}

// GasFeeCapCmp compares the fee cap of two transactions.
func (tx *Transaction) GasFeeCapCmp(other *Transaction) int {
	return tx.inner.gasFeeCap().Cmp(other.inner.gasFeeCap())
}

// GasFeeCapIntCmp compares the fee cap of the transaction against the given fee cap.
func (tx *Transaction) GasFeeCapIntCmp(other *big.Int) int {
	return tx.inner.gasFeeCap().Cmp(other)
}

// GasTipCapCmp compares the gasTipCap of two transactions.
func (tx *Transaction) GasTipCapCmp(other *Transaction) int {
	return tx.inner.gasTipCap().Cmp(other.inner.gasTipCap())
}

// GasTipCapIntCmp compares the gasTipCap of the transaction against the given gasTipCap.
func (tx *Transaction) GasTipCapIntCmp(other *big.Int) int {
	return tx.inner.gasTipCap().Cmp(other)
}

// EffectiveGasTip returns the effective miner gasTipCap for the given base fee.
// Note: if the effective gasTipCap is negative, this method returns both error
// the actual negative value, _and_ ErrGasFeeCapTooLow
func (tx *Transaction) EffectiveGasTip(baseFee *big.Int) (*big.Int, error) {
	if baseFee == nil {
		return tx.GasTipCap(), nil
	}
	var err error
	gasFeeCap := tx.GasFeeCap()
	if gasFeeCap.Cmp(baseFee) < 0 {
		err = ErrGasFeeCapTooLow
	}
	gasFeeCap = gasFeeCap.Sub(gasFeeCap, baseFee)

	gasTipCap := tx.GasTipCap()
	if gasTipCap.Cmp(gasFeeCap) < 0 {
		return gasTipCap, err
	}
	return gasFeeCap, err
}

// EffectiveGasTipValue is identical to EffectiveGasTip, but does not return an
// error in case the effective gasTipCap is negative
func (tx *Transaction) EffectiveGasTipValue(baseFee *big.Int) *big.Int {
	effectiveTip, _ := tx.EffectiveGasTip(baseFee)
	return effectiveTip
}

// EffectiveGasTipCmp compares the effective gasTipCap of two transactions assuming the given base fee.
func (tx *Transaction) EffectiveGasTipCmp(other *Transaction, baseFee *big.Int) int {
	if baseFee == nil {
		return tx.GasTipCapCmp(other)
	}
	return tx.EffectiveGasTipValue(baseFee).Cmp(other.EffectiveGasTipValue(baseFee))
}

// EffectiveGasTipIntCmp compares the effective gasTipCap of a transaction to the given gasTipCap.
func (tx *Transaction) EffectiveGasTipIntCmp(other *big.Int, baseFee *big.Int) int {
	if baseFee == nil {
		return tx.GasTipCapIntCmp(other)
	}
	return tx.EffectiveGasTipValue(baseFee).Cmp(other)
}

// BlobGas returns the blob gas limit of the transaction for blob transactions, 0 otherwise.
func (tx *Transaction) BlobGas() uint64 {
	if blobtx, ok := tx.inner.(*BlobTx); ok {
		return blobtx.blobGas()
	}
	return 0
}

// BlobGasFeeCap returns the blob gas fee cap per blob gas of the transaction for blob transactions, nil otherwise.
func (tx *Transaction) BlobGasFeeCap() *big.Int {
	if blobtx, ok := tx.inner.(*BlobTx); ok {
		return blobtx.BlobFeeCap.ToBig()
	}
	return nil
}

// BlobHashes returns the hashes of the blob commitments for blob transactions, nil otherwise.
func (tx *Transaction) BlobHashes() []common.Hash {
	if blobtx, ok := tx.inner.(*BlobTx); ok {
		return blobtx.BlobHashes
	}
	return nil
}

// BlobTxSidecar returns the sidecar of a blob transaction, nil otherwise.
func (tx *Transaction) BlobTxSidecar() *BlobTxSidecar {
	if blobtx, ok := tx.inner.(*BlobTx); ok {
		return blobtx.Sidecar
	}
	return nil
}

// BlobGasFeeCapCmp compares the blob fee cap of two transactions.
func (tx *Transaction) BlobGasFeeCapCmp(other *Transaction) int {
	return tx.BlobGasFeeCap().Cmp(other.BlobGasFeeCap())
}

// BlobGasFeeCapIntCmp compares the blob fee cap of the transaction against the given blob fee cap.
func (tx *Transaction) BlobGasFeeCapIntCmp(other *big.Int) int {
	return tx.BlobGasFeeCap().Cmp(other)
}

// WithoutBlobTxSidecar returns a copy of tx with the blob sidecar removed.
func (tx *Transaction) WithoutBlobTxSidecar() *Transaction {
	blobtx, ok := tx.inner.(*BlobTx)
	if !ok {
		return tx
	}
	cpy := &Transaction{
		inner: blobtx.withoutSidecar(),
		time:  tx.time,
	}
	// Note: tx.size cache not carried over because the sidecar is included in size!
	if h := tx.hash.Load(); h != nil {
		cpy.hash.Store(h)
	}
	if f := tx.from.Load(); f != nil {
		cpy.from.Store(f)
	}
	return cpy
}

// WithBlobTxSidecar returns a copy of tx with the blob sidecar added.
func (tx *Transaction) WithBlobTxSidecar(sideCar *BlobTxSidecar) *Transaction {
	blobtx, ok := tx.inner.(*BlobTx)
	if !ok {
		return tx
	}
	cpy := &Transaction{
		inner: blobtx.withSidecar(sideCar),
		time:  tx.time,
	}
	// Note: tx.size cache not carried over because the sidecar is included in size!
	if h := tx.hash.Load(); h != nil {
		cpy.hash.Store(h)
	}
	if f := tx.from.Load(); f != nil {
		cpy.from.Store(f)
	}
	return cpy
}

// SetCodeAuthorizations returns the authorizations list of the transaction.
func (tx *Transaction) SetCodeAuthorizations() []SetCodeAuthorization {
	setcodetx, ok := tx.inner.(*SetCodeTx)
	if !ok {
		return nil
	}
	return setcodetx.AuthList
}

// SetCodeAuthorities returns a list of unique authorities from the
// authorization list.
func (tx *Transaction) SetCodeAuthorities() []common.Address {
	setcodetx, ok := tx.inner.(*SetCodeTx)
	if !ok {
		return nil
	}
	var (
		marks = make(map[common.Address]bool)
		auths = make([]common.Address, 0, len(setcodetx.AuthList))
	)
	for _, auth := range setcodetx.AuthList {
		if addr, err := auth.Authority(); err == nil {
			if marks[addr] {
				continue
			}
			marks[addr] = true
			auths = append(auths, addr)
		}
	}
	return auths
}

// SetTime sets the decoding time of a transaction. This is used by tests to set
// arbitrary times and by persistent transaction pools when loading old txs from
// disk.
func (tx *Transaction) SetTime(t time.Time) {
	tx.time = t
}

// Time returns the time when the transaction was first seen on the network. It
// is a heuristic to prefer mining older txs vs new all other things equal.
func (tx *Transaction) Time() time.Time {
	return tx.time
}

// Hash returns the transaction hash.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return *hash
	}

	var h common.Hash
	if tx.Type() == LegacyTxType {
		h = rlpHash(tx.inner)
	} else {
		h = prefixedRlpHash(tx.Type(), tx.inner)
	}
	tx.hash.Store(&h)
	return h
}

// Size returns the true encoded storage size of the transaction, either by encoding
// and returning it, or returning a previously cached value.
func (tx *Transaction) Size() uint64 {
	if size := tx.size.Load(); size > 0 {
		return size
	}

	// Cache miss, encode and cache.
	// Note we rely on the assumption that all tx.inner values are RLP-encoded!
	c := writeCounter(0)
	rlp.Encode(&c, &tx.inner)
	size := uint64(c)

	// For blob transactions, add the size of the blob content and the outer list of the
	// tx + sidecar encoding.
	if sc := tx.BlobTxSidecar(); sc != nil {
		size += rlp.ListSize(sc.encodedSize())
	}

	// For typed transactions, the encoding also includes the leading type byte.
	if tx.Type() != LegacyTxType {
		size += 1
	}

	tx.size.Store(size)
	return size
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	if r == nil || s == nil || v == nil {
		return nil, fmt.Errorf("%w: r: %s, s: %s, v: %s", ErrInvalidSig, r, s, v)
	}
	cpy := tx.inner.copy()
	cpy.setSignatureValues(signer.ChainID(), v, r, s)
	return &Transaction{inner: cpy, time: tx.time}, nil
}

// Transactions implements DerivableList for transactions.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// EncodeIndex encodes the i'th transaction to w. Note that this does not check for errors
// because we assume that *Transaction will only ever contain valid txs that were either
// constructed by decoding or via public API in this package.
func (s Transactions) EncodeIndex(i int, w *bytes.Buffer) {
	tx := s[i]
	if tx.Type() == LegacyTxType {
		rlp.Encode(w, tx.inner)
	} else {
		tx.encodeTyped(w)
	}
}

// TxDifference returns a new set of transactions that are present in a but not in b.
func TxDifference(a, b Transactions) Transactions {
	keep := make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{}, b.Len())
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// HashDifference returns a new set of hashes that are present in a but not in b.
func HashDifference(a, b []common.Hash) []common.Hash {
	keep := make([]common.Hash, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, hash := range b {
		remove[hash] = struct{}{}
	}

	for _, hash := range a {
		if _, ok := remove[hash]; !ok {
			keep = append(keep, hash)
		}
	}

	return keep
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int           { return len(s) }
func (s TxByNonce) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s TxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// copyAddressPtr copies an address.
func copyAddressPtr(a *common.Address) *common.Address {
	if a == nil {
		return nil
	}
	cpy := *a
	return &cpy
}

```

