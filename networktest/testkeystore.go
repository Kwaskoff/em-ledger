// This software is Copyright (c) 2019-2020 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package networktest

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
)

const (
	KeyPwd   = "pwd12345"
	Bip39Pwd = ""
)

type (
	KeyStore struct {
		path    string
		keybase keys.Keybase

		Authority,
		DeputyKey,
		Key1,
		Key2,
		Key3,
		Key4,
		Key5,
		Key6 Key

		MultiKey Key

		Validators []Key
	}

	Key struct {
		name    string
		keybase keys.Keybase
		privkey crypto.PrivKey
		pubkey  crypto.PubKey
		address sdk.AccAddress
	}
)

func newKey(name string, keybase keys.Keybase) Key {
	// Extract key information to prevent future keystore access. Makes concurrent key usage possible.
	var (
		privkey, _ = keybase.ExportPrivateKeyObject(name, KeyPwd)
		info, _    = keybase.Get(name)
	)

	var address sdk.AccAddress
	var pubKey crypto.PubKey
	if info != nil {
		pubKey = info.GetPubKey()
		address = info.GetAddress()
	}

	return Key{
		name:    name,
		keybase: keybase,
		privkey: privkey,
		pubkey:  pubKey,
		address: address,
	}
}

func (k Key) GetAddress() string {
	if k.address.Empty() {
		info, err := k.keybase.Get(k.name)
		if err != nil {
			panic(err)
		}

		k.address = info.GetAddress()
	}

	return k.address.String()
}

func (k Key) GetPublicKey() crypto.PubKey {
	if k.pubkey != nil {
		return k.pubkey
	}

	return k.privkey.PubKey()
}

func (k Key) Sign(bz []byte) ([]byte, error) {
	return k.privkey.Sign(bz)
}

func NewKeystore() (*KeyStore, error) {
	path, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	keybase, err := keys.NewKeyring(sdk.KeyringServiceName(), keys.BackendTest, path, nil)
	if err != nil {
		return nil, err
	}

	initializeKeystore(keybase)

	// TODO This looks kind of horrible. Refactor to something prettier
	ks := &KeyStore{
		keybase:   keybase,
		path:      path,
		Authority: newKey("authoritykey", keybase),
		DeputyKey: newKey("deputykey", keybase),
		Key1:      newKey("key1", keybase),
		Key2:      newKey("key2", keybase),
		Key3:      newKey("key3", keybase),
		Key4:      newKey("key3", keybase),
		Key5:      newKey("key3", keybase),
		Key6:      newKey("key3", keybase),

		MultiKey: newKey("multikey", keybase),
		Validators: []Key{
			newKey("validator0", keybase),
			newKey("validator1", keybase),
			newKey("validator2", keybase),
			newKey("validator3", keybase),
		},
	}

	return ks, nil
}

func (ks KeyStore) Close() {
	_ = os.RemoveAll(ks.path)
}

func (ks KeyStore) GetPath() string {
	return ks.path
}

func (ks KeyStore) String() string {
	keyinfos, err := ks.keybase.List()
	if err != nil {
		return err.Error()
	}

	var sb strings.Builder
	for _, info := range keyinfos {
		sb.WriteString(fmt.Sprintf("%v - %v (%v)\n", info.GetName(), info.GetAddress().String(), info.GetAlgo()))
	}

	return sb.String()
}

func (ks KeyStore) addValidatorKeys(testnetoutput string) {
	scan := bufio.NewScanner(strings.NewReader(testnetoutput))
	seeds := make([]string, 0)
	for scan.Scan() {
		s := scan.Text()
		if strings.Contains(s, "Key mnemonic for Validator") {
			seed := strings.Split(s, ":")[1]
			seeds = append(seeds, strings.TrimSpace(seed))
		}
	}

	for i, mnemonic := range seeds {
		accountName := fmt.Sprintf("validator%v", i)
		hdPath := sdk.GetConfig().GetFullFundraiserPath()
		_, err := ks.keybase.CreateAccount(accountName, mnemonic, Bip39Pwd, KeyPwd, hdPath, keys.Secp256k1)
		if err != nil {
			panic(err)
		}
	}
}

func (ks KeyStore) addDeputyKey() {
	mn := "play witness auto coast domain win tiny dress glare bamboo rent mule delay exact arctic vacuum laptop hidden siren sudden six tired fragile penalty"
	// create the deputy account
	hdPath := sdk.GetConfig().GetFullFundraiserPath()
	deputyAccount, err := ks.keybase.CreateAccount("deputykey", mn, "", KeyPwd, hdPath,
		keys.Secp256k1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("deputy address: %s\nmnemonic: %s\n", deputyAccount.GetAddress().String(), mn)
}

func initializeKeystore(kb keys.Keybase) {
	hdPath := sdk.GetConfig().GetFullFundraiserPath()
	const mnemonic1 = "then nuclear favorite advance plate glare shallow enhance replace embody list dose quick scale service sentence hover announce advance nephew phrase order useful this"
	ac1, _ := kb.CreateAccount("authoritykey", mnemonic1, "", KeyPwd, hdPath, keys.Secp256k1)
	fmt.Printf("Created account %s from mnemonic: %s\n", ac1.GetAddress(), mnemonic1)

	const mnemonic2 = "document weekend believe whip diesel earth hope elder quiz pact assist quarter public deal height pulp roof organ animal health month holiday front pencil"
	ac2, _ := kb.CreateAccount("key1", mnemonic2, "", KeyPwd, hdPath, keys.Secp256k1)
	fmt.Printf("Created account %s from mnemonic: %s\n", ac2.GetAddress(), mnemonic2)

	const mnemonic3 = "treat ocean valid motor life marble syrup lady nephew grain cherry remember lion boil flock outside cupboard column dad rare build nut hip ostrich"
	ac3, _ := kb.CreateAccount("key2", mnemonic3, "", KeyPwd, hdPath, keys.Secp256k1)
	fmt.Printf("Created account %s from mnemonic: %s\n", ac3.GetAddress(), mnemonic3)

	const mnemonic4 = "rice short length buddy zero snake picture enough steak admit balance garage exit crazy cloud this sweet virus can aunt embrace picnic stick wheel"
	ac4, _ := kb.CreateAccount("key3", mnemonic4, "", KeyPwd, hdPath, keys.Secp256k1)
	fmt.Printf("Created account %s from mnemonic: %s\n", ac4.GetAddress(), mnemonic4)

	const mnemonic5 = "census museum crew rude tower vapor mule rib weasel faith page cushion rain inherit much cram that blanket occur region track hub zero topple"
	ac5, _ := kb.CreateAccount("key4", mnemonic5, "", KeyPwd, hdPath, keys.Secp256k1)
	fmt.Printf("Created account %s from mnemonic: %s\n", ac5.GetAddress(), mnemonic5)

	const mnemonic6 = "flavor print loyal canyon expand salmon century field say frequent human dinosaur frame claim bridge affair web way direct win become merry crash frequent"
	ac6, _ := kb.CreateAccount("key5", mnemonic6, "", KeyPwd, hdPath, keys.Secp256k1)
	fmt.Printf("Created account %s from mnemonic: %s\n", ac6.GetAddress(), mnemonic6)

	const mnemonic7 = "very health column only surface project output absent outdoor siren reject era legend legal twelve setup roast lion rare tunnel devote style random food"
	ac7, _ := kb.CreateAccount("key6", mnemonic7, "", KeyPwd, hdPath, keys.Secp256k1)
	fmt.Printf("Created account %s from mnemonic: %s\n", ac7.GetAddress(), mnemonic7)

	// Create a multisig key entry consisting of key1, key2 and key3 with a threshold of 2
	pks := make([]crypto.PubKey, 3)
	for i, keyname := range []string{"key1", "key2", "key3"} {
		keyinfo, err := kb.Get(keyname)
		if err != nil {
			panic(err)
		}

		pks[i] = keyinfo.GetPubKey()
	}

	pk := multisig.NewPubKeyMultisigThreshold(2, pks)
	_, err := kb.CreateMulti("multikey", pk)
	if err != nil {
		panic(err)
	}
}
