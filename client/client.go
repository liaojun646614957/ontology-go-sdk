/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/types"
	bc "github.com/ontio/ontology/http/base/common"
	"github.com/tendermint/iavl"
)

type ClientMgr struct {
	rpc       *RpcClient  //Rpc client used the rpc api of ontology
	rest      *RestClient //Rest client used the rest api of ontology
	ws        *WSClient   //Web socket client used the web socket api of ontology
	defClient OntologyClient
	qid       uint64
}

type Layer2ClientMgr struct {
	client *ClientMgr
}

func NewLayer2ClientMgr(c *ClientMgr) *Layer2ClientMgr {
	layer2Client := &Layer2ClientMgr{
		client: c,
	}
	return layer2Client
}

func (this *ClientMgr) NewRpcClient() *RpcClient {
	this.rpc = NewRpcClient()
	return this.rpc
}

func (this *ClientMgr) GetRpcClient() *RpcClient {
	return this.rpc
}

func (this *ClientMgr) NewRestClient() *RestClient {
	this.rest = NewRestClient()
	return this.rest
}

func (this *ClientMgr) GetRestClient() *RestClient {
	return this.rest
}

func (this *ClientMgr) NewWebSocketClient() *WSClient {
	wsClient := NewWSClient()
	this.ws = wsClient
	return wsClient
}

func (this *ClientMgr) GetWebSocketClient() *WSClient {
	return this.ws
}

func (this *ClientMgr) SetDefaultClient(client OntologyClient) {
	this.defClient = client
}

func (this *ClientMgr) GetCurrentBlockHeight() (uint32, error) {
	client := this.getClient()
	if client == nil {
		return 0, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getCurrentBlockHeight(this.getNextQid())
	if err != nil {
		return 0, err
	}
	return utils.GetUint32(data)
}

func (this *ClientMgr) GetCurrentBlockHash() (common.Uint256, error) {
	client := this.getClient()
	if client == nil {
		return common.UINT256_EMPTY, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getCurrentBlockHash(this.getNextQid())
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return utils.GetUint256(data)
}

func (this *ClientMgr) GetBlockByHeight(height uint32) (*types.Block, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockByHeight(this.getNextQid(), height)
	if err != nil {
		return nil, err
	}
	return utils.GetBlock(data)
}

func (this *Layer2ClientMgr) GetLayer2BlockByHeight(height uint32) (*utils.Layer2Block, error) {
	client := this.client.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockByHeight(this.client.getNextQid(), height)
	if err != nil {
		return nil, err
	}
	return utils.GetLayer2Block(data)
}

func (this *ClientMgr) GetBlockInfoByHeight(height uint32) ([]byte, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockInfoByHeight(this.getNextQid(), height)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *ClientMgr) GetBlockByHash(blockHash string) (*types.Block, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockByHash(this.getNextQid(), blockHash)
	if err != nil {
		return nil, err
	}
	return utils.GetBlock(data)
}

func (this *Layer2ClientMgr) GetLayer2BlockByHash(blockHash string) (*utils.Layer2Block, error) {
	client := this.client.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockByHash(this.client.getNextQid(), blockHash)
	if err != nil {
		return nil, err
	}
	return utils.GetLayer2Block(data)
}

func (this *ClientMgr) GetTransaction(txHash string) (*types.Transaction, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getRawTransaction(this.getNextQid(), txHash)
	if err != nil {
		return nil, err
	}
	return utils.GetTransaction(data)
}

func (this *ClientMgr) GetBlockHash(height uint32) (common.Uint256, error) {
	client := this.getClient()
	if client == nil {
		return common.UINT256_EMPTY, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockHash(this.getNextQid(), height)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return utils.GetUint256(data)
}

func (this *ClientMgr) GetBlockHeightByTxHash(txHash string) (uint32, error) {
	client := this.getClient()
	if client == nil {
		return 0, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockHeightByTxHash(this.getNextQid(), txHash)
	if err != nil {
		return 0, err
	}
	return utils.GetUint32(data)
}

func (this *ClientMgr) GetBlockTxHashesByHeight(height uint32) (*sdkcom.BlockTxHashes, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getBlockTxHashesByHeight(this.getNextQid(), height)
	if err != nil {
		return nil, err
	}
	return utils.GetBlockTxHashes(data)
}

func (this *ClientMgr) GetStorage(contractAddress string, key []byte) ([]byte, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getStorage(this.getNextQid(), contractAddress, key)
	if err != nil {
		return nil, err
	}
	return utils.GetStorage(data)
}

func (this *ClientMgr) GetSmartContract(contractAddress string) (*payload.DeployCode, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getSmartContract(this.getNextQid(), contractAddress)
	if err != nil {
		return nil, err
	}
	return utils.GetSmartContract(data)
}

func (this *ClientMgr) GetSmartContractEvent(txHash string) (*sdkcom.SmartContactEvent, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getSmartContractEvent(this.getNextQid(), txHash)
	if err != nil {
		return nil, err
	}
	return utils.GetSmartContractEvent(data)
}

func (this *ClientMgr) GetSmartContractEventByBlock(height uint32) ([]*sdkcom.SmartContactEvent, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getSmartContractEventByBlock(this.getNextQid(), height)
	if err != nil {
		return nil, err
	}
	if data == nil || string(data) == "" || string(data) == "\"\"" {
		return nil, nil
	}
	return utils.GetSmartContactEvents(data)
}

func (this *ClientMgr) GetMerkleProof(txHash string) (*sdkcom.MerkleProof, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getMerkleProof(this.getNextQid(), txHash)
	if err != nil {
		return nil, err
	}
	return utils.GetMerkleProof(data)
}

func (this *ClientMgr) GetCrossStatesProof(height uint32, key []byte) (*bc.CrossStatesProof, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getCrossStatesProof(this.getNextQid(), height, key)
	if err != nil {
		return nil, err
	}
	return utils.GetCrossStatesProof(data)
}

func (this *ClientMgr) GetCrossChainMsg(height uint32) (string, error) {
	client := this.getClient()
	if client == nil {
		return "", fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getCrossChainMsg(this.getNextQid(), height)
	if err != nil {
		return "", err
	}
	return utils.GetCrossChainMsg(data)
}

func (this *ClientMgr) GetMemPoolTxState(txHash string) (*sdkcom.MemPoolTxState, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getMemPoolTxState(this.getNextQid(), txHash)
	if err != nil {
		return nil, err
	}
	return utils.GetMemPoolTxState(data)
}

func (this *ClientMgr) GetMemPoolTxCount() (*sdkcom.MemPoolTxCount, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getMemPoolTxCount(this.getNextQid())
	if err != nil {
		return nil, err
	}
	return utils.GetMemPoolTxCount(data)
}

func (this *ClientMgr) GetVersion() (string, error) {
	client := this.getClient()
	if client == nil {
		return "", fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getVersion(this.getNextQid())
	if err != nil {
		return "", err
	}
	return utils.GetVersion(data)
}

func (this *ClientMgr) GetNetworkId() (uint32, error) {
	client := this.getClient()
	if client == nil {
		return 0, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getNetworkId(this.getNextQid())
	if err != nil {
		return 0, err
	}
	return utils.GetUint32(data)
}

func (this *ClientMgr) SendTransaction(mutTx *types.MutableTransaction) (common.Uint256, error) {
	client := this.getClient()
	if client == nil {
		return common.UINT256_EMPTY, fmt.Errorf("don't have available client of ontology")
	}
	tx, err := mutTx.IntoImmutable()
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	data, err := client.sendRawTransaction(this.getNextQid(), tx, false)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return utils.GetUint256(data)
}

func (this *ClientMgr) PreExecTransaction(mutTx *types.MutableTransaction) (*sdkcom.PreExecResult, error) {
	client := this.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	tx, err := mutTx.IntoImmutable()
	if err != nil {
		return nil, err
	}
	data, err := client.sendRawTransaction(this.getNextQid(), tx, true)
	if err != nil {
		return nil, err
	}
	preResult := &sdkcom.PreExecResult{}
	err = json.Unmarshal(data, &preResult)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal PreExecResult:%s error:%s", data, err)
	}
	return preResult, nil
}

func (this *Layer2ClientMgr) VerifyLayer2StoreProof(key []byte, value []byte, proof []byte, stateRoot []byte) (bool, error) {
	source := common.NewZeroCopySource(proof)
	storeProof := new(utils.Layer2StoreProof)
	err := storeProof.Deserialization(source)
	if err != nil {
		return false, err
	}
	proof_iavl := iavl.RangeProof(*storeProof)
	err = proof_iavl.Verify(stateRoot)
	if err != nil {
		return false, err
	}
	err = proof_iavl.VerifyItem(key, value)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *Layer2ClientMgr) GetLayer2StoreKey(contractAddress string, key []byte) ([]byte, error) {
	newKey := make([]byte, 0)
	if len(contractAddress) > 0 {
		newKey = append(newKey, byte(0x05))
		contractAddrBytes, _ := hex.DecodeString(contractAddress)
		contractAddr := common.ToArrayReverse(contractAddrBytes)
		newKey = append(newKey, contractAddr...)
		newKey = append(newKey, key...)
	} else {
		newKey = append(newKey, key...)
	}
	return newKey, nil
}

func (this *Layer2ClientMgr) GetLayer2StoreProof(key []byte) (*sdkcom.Layer2StoreProof, error) {
	client := this.client.getClient()
	if client == nil {
		return nil, fmt.Errorf("don't have available client of ontology")
	}
	data, err := client.getLayer2StoreProof(this.client.getNextQid(), key)
	if err != nil {
		return nil, err
	}
	return utils.GetLayer2StoreProof(data)
}

//WaitForGenerateBlock Wait ontology generate block. Default wait 2 blocks.
//return timeout error when there is no block generate in some time.
func (this *ClientMgr) WaitForGenerateBlock(timeout time.Duration, blockCount ...uint32) (bool, error) {
	count := uint32(2)
	if len(blockCount) > 0 && blockCount[0] > 0 {
		count = blockCount[0]
	}
	blockHeight, err := this.GetCurrentBlockHeight()
	if err != nil {
		return false, fmt.Errorf("GetCurrentBlockHeight error:%s", err)
	}
	secs := int(timeout / time.Second)
	if secs <= 0 {
		secs = 1
	}
	for i := 0; i < secs; i++ {
		time.Sleep(time.Second)
		curBlockHeigh, err := this.GetCurrentBlockHeight()
		if err != nil {
			continue
		}
		if curBlockHeigh-blockHeight >= count {
			return true, nil
		}
	}
	return false, fmt.Errorf("timeout after %d (s)", secs)
}

func (this *ClientMgr) getClient() OntologyClient {
	if this.defClient != nil {
		return this.defClient
	}
	if this.rpc != nil {
		return this.rpc
	}
	if this.rest != nil {
		return this.rest
	}
	if this.ws != nil {
		return this.ws
	}
	return nil
}

func (this *ClientMgr) getNextQid() string {
	return fmt.Sprintf("%d", atomic.AddUint64(&this.qid, 1))
}
