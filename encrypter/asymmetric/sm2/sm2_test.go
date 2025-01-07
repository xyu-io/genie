package sm2_lib

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"log"
	"os"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	// 生成 SM2 椭圆曲线
	privateKey, _ := sm2.GenerateKey(rand.Reader)

	// 获取公钥
	pubKey := privateKey.PublicKey

	// 将公钥转成字符串
	pubKeyStr := hex.EncodeToString(sm2.Compress(&pubKey))

	fmt.Println("SM2 public key:", pubKeyStr)
	// 将私钥转换为 16 进制字符串
	//2.通过x509将私钥反序列化并进行pem编码
	privateKeyToPem, err := x509.WritePrivateKeyToPem(privateKey, nil)
	if err != nil {
		panic(err)
	}

	pubKeyToPem, err := x509.WritePublicKeyToPem(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("pub.pem")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(pubKeyToPem)
	if err != nil {
		return
	}
	// 将私钥保存到文件
	//3.将私钥写入磁盘文件
	file, err = os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(privateKeyToPem)
	if err != nil {
		return
	}
}

func TestEncrypt(t *testing.T) {
	var data = map[string]string{
		"data": "hello",
	}

	// 生成 SM2 椭圆曲线
	privateKey, _ := sm2.GenerateKey(rand.Reader)

	// 获取公钥
	pubKey := privateKey.PublicKey

	marshal, _ := json.Marshal(data)

	cipherBytes, err := sm2.Encrypt(&pubKey, marshal, nil, sm2.C1C3C2)
	if err != nil {
		log.Println(err)
		return
	}
	ciphertext := hex.EncodeToString(cipherBytes)

	fmt.Println(ciphertext)
}

func TestDecrypt(t *testing.T) {
	var ciphertext = "0429a4f002bfb19ff1d19aa839b7c15f5951c3306388e0748aff9b7337a265c23869bfcd5648c580641c49494ebf964775cade6187d188638e6a6ac239c5b134d4b235fae0eebd9a8b52455ac7c219701f947ee77d6c912f869ac8802e890ef235ee6916b860c8d74e58f52130476c0ed5e22a9dbcf56d764b252ccb5392a0e4e4b4cb2032c49cb6a2af73183f72a3bd65ee805b68676ce0a47f216c057fc24540eed857648d85c0277985d4c650f2f25a3cffe3acf7fce67f52305b14c6981c175b57e6d3bbdd95b882066762246976d4a460cde5e6dacba8c51039bcb62d5ae39b3d49ba10373f412abf1fa2e47341db941ddd5c4c82a0e197dba66170909d84cb3d9064e4300f5cdd872579a61aa0d571b25579f76a83b7e1f7f9e51a05fc2e3d99c9f6d158b7046787dd8f8f37292812757c660a792509b989e4ae29e593659cf0e0c14ae559f589a38d7d89a736f388122caba095c76f72b2778115da4aea7183cf491d91f47239351fc1196d3d58b0e19ace40128c59dfe6c4a3adeec6f5553b856302d07d73e36f1dcab71f3824badbc48f27b5a129f97ad74f319ddaed3329f5cdebc45f03b0c591765848c7191a6745207010442f5f32401a3554640716f33e8507a2372139b7e0bcc55d543595d5c3388aab32d32eff3c3eeb1e95c95d078a8e8c42afc93d29da2356a3a0490b804a9bac98ed4157a466804ef018124e1467df1771de96ddbe74c2874d865aa784b522fdc8255b8262c756880368214b839ed2dc094a2283267eb4c0d042089602ffe6124c93bf0a8fd4d88c76325ccbfade7a7ba399acd604d7111e10f32f985f5986ba2b2a72e6748baab89aca488d04c10619add57261037c0c5ca3318a527f501138a85edcc8be7d696562408b72306ce54d3990096e35b8330d379380546db1d3384b2be64ffe9441ad71f7cd03d181da307d0f60c597a3f2b151b0f9ab3a4b4aefcaf8c40a0a4f06a4b191ea1d9d65040a661f01c3385537ba7df28718468cf6b4c774a719fabc0a4148c7487732e698c26eddae6a205145e27d2ce291cff59e5f42c3a0fde5cdfe05b7def1981f0e1b5a4235de59715358211b893450e2e0e7734338ab7e03aee20a283346cc59ddce8e55998f48cc9717245e85ff81dc3f65ea916b51514a47e7147ee5833091bdc5c31c4604d536beb44449e22500b986fdba072e969821973bb243f8d5293f36eda8592a18075cc6486da9798942d578c1988ca39f0790ea14a858401e2ad720d875e7fea939012ba154e9571d8007d271b37374a353e5d6bb02b036f353af50ecd824db0bedfa26bc8db0f9f0f12f90717b17dcc5c234aabeda907e53f87bb355a660d59e131470fdc23fa8a56de1a81d809f60a29acabdebb988e7a391eaf8b63f3589491945e98001b627d008e02e9a5da6580640de6acda07ec8bba3f96218634634fd9fffe810cd651b06a372dcb9f121171f08f175f2f70a0359664f440a6c1534e5f046b034f38b5dda7d108a4499d26634ddb868c225ce3ad66d80fac29e205506898199dabc67b6ff56f81b79b5460939bc0ad7e310a7626d2150f565db4a735fd61b322a237ff9c74efff082a4c1aa09fd8c93c4f85c550923f0e3a361f30097f533a44180a6e925597e38d00307a68e4fc8357b4d783738ff7c7b2f7f923dbfce83f32e566291249976cd970ce8d356c9b645b9fe1924d6339ad3df831e3878717fc3b70e9d7a71c5625fa252f589f6a3a7ed0585c467575f2cfda55fe02497f692c32794f8227844953f99c29a9cfa698d4f6c0d0f4e281f508091d40822277f3038cf360962fe0a04d6753016bb493aea92bfc738cfa2ccfa0faa178e82bca59d22aae66d49a3170b006aed708ba1b401f32fee3c0a0b9215bddb692caf26f993475da6178fe8c695925ba6c0155ee065294e6e6f8394d3dee90774fc289507d3ba1348ea631a37bfe549ac4ba8f88f123032bd4f369778ef9507e057227b9f37ce6221c0c6767f3505aad83fa37d17516a6fb2dc9a1043333e42db0f12af6811f0127c1e22120620b41a878ca8f2724ac2730fd21c1510cf14e56ee3e3127111334a82b5eb4f6e8afae24b38e0de551585a7488f3915e3b29a9021782774cf398e32e78747c3caa0e097cd03f5f92850013ee776782f661416d56236d52c75d26510de39428ca3dad68dafc226d049f4e8cbc06288cf9e2a6f1471456ecf1eba6e880a924c7716aa149d2ada570864c1d7a36cdc8c186fc9b0b771f6466c970b1ce7626d70523bb1fe51861e399acb09f3a8c21c190e41e30c4f10922ac96cb54b369ec84d5a135a17d4ad786e850dc5a92ce616800a98f293351d1d1968cb01592f83eff57dc8435e688cf6cb5418a026197cfe283de955d05a349d6038c0740eb0796fd9fbe9b9b5648037d3d01654ec7ea96e6d6b2f5e17c79c3fd8e8f944985618d592ba7ef2d90ca17c6ca8e72ba4eec0fbff63da1d6fd66643029acafca181bd4b74a90a971c2ed50a07294b63b010c89492514615ffe72d3b55767f7a7ed50e05e5754aef4bfba30a094d4af938258d598682c26a"
	file, err := os.Open("private.pem")
	if err != nil {
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)

	//2.将pem格式私钥文件解码并反序列话
	privateKeyFromPem, err := x509.ReadPrivateKeyFromPem(buf, nil)
	//3.解密
	cipherByte, err := hex.DecodeString(ciphertext)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	planiText, err := sm2.Decrypt(privateKeyFromPem, cipherByte, sm2.C1C3C2)
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println(string(planiText))
}
