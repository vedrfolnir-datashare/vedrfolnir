package main

import (
	core "Ved/core"
	"Ved/lib/VedCrypto"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	//"strconv"
	"strings"
	"unsafe"
)

var HSMT = make(map[int]*core.MerkleNode)
var Data = make(map[int]int)
var DataAPatient []string
var DataAllPatient = make(map[int][]string)
var key0 []byte = []byte("ThisIs16ByteKey0")
var key1 []byte = []byte("ThisIs16ByteKey1")

func main_temp1() {
	//TestIndex()
	//TestIndexInit()
	TestIndexConstruction()
	//TestIndexVerify()
	//TestIndexVerifyBandwidth()
	//TestIndexProcess()
	//TestIndexVerifyTimeSpace()
	TestIndexVerifyTimeSpaceAll()
}

func TestIndexProcess() {
	kw := "apple"

	time0 := time.Now()
	HSMT[0] = core.InitIndex(key0)
	time1 := time.Now()
	kwHashStr := string(VedCrypto.Hash([]byte(kw)))
	lastChar := []byte{kwHashStr[len(kwHashStr)-1]}
	tempProof, _ := VedCrypto.Encrypt(lastChar, key0)
	nonexistenceProof := VedCrypto.Hash(tempProof, tempProof)
	time2 := time.Now()

	proofs, _, _ := core.CheckExistence([]byte(kw), HSMT[0], nonexistenceProof)
	time3 := time.Now()

	BHT_CMT := core.NewIndex(key1, 1)
	time4 := time.Now()

	HSMT[0] = core.UpdateIndex(proofs, BHT_CMT)
	time5 := time.Now()

	initTime := time1.Sub(time0)
	keyGenerateTime := time2.Sub(time1)
	proofCheckTime := time3.Sub(time2)
	indexGenerateTime := time4.Sub(time3)
	indexUpdateTime := time5.Sub(time4)

	fmt.Printf("Index init %s, key generate %s, proof check %s, BHT-CMT generate %s, and update: %s \n", initTime, keyGenerateTime, proofCheckTime, indexGenerateTime, indexUpdateTime)
}

func TestIndex() {
	totalTime := time.Duration(0)

	time0 := time.Now()

	HSMT[0] = core.InitIndex(key0)

	time1 := time.Now()
	data := "this is a medical data, which contains many many hard core medical worlds. this is a medical data, which contains many many hard core medical worlds. this is a medical data, which contains many many hard core medical worlds. this is a medical data, which contains many many hard core medical worlds. this is a medical data, which contains"
	keyword := strings.Split(data, " ")

	time2 := time.Now()
	for _, kw := range keyword {
		// 检查关键字是否存在，不存在添加关键字

		// 这一段复杂的字符串转换操作是构建不存在证明
		kwHashStr := string(VedCrypto.Hash([]byte(kw)))
		lastChar := []byte{kwHashStr[len(kwHashStr)-1]}
		tempProof, _ := VedCrypto.Encrypt(lastChar, key0)
		nonexistenceProof := VedCrypto.Hash(tempProof, tempProof)

		proofs, res, _ := core.CheckExistence([]byte(kw), HSMT[0], nonexistenceProof)
		if res == false {
			//print("keyword ", kw, " doest not exist, create an BHT_CMT \n")
			//time2_1 := time.Now()
			BHT_CMT := core.NewIndex(key1, 1)
			time2_2 := time.Now()
			HSMT[0] = core.UpdateIndex(proofs, BHT_CMT)
			time2_3 := time.Now()
			totalTime += time2_3.Sub(time2_2)
		}
	}

	time3 := time.Now()

	initTime := time1.Sub(time0)
	updateTime := time3.Sub(time2)

	totalTime = totalTime / time.Duration(len(keyword))

	fmt.Printf("One record Index init %s, and update: %s. Each key word average: %s \n", initTime, updateTime, totalTime)
}

func TestIndexConstruction() {
	db, err := sqlx.Open("mysql", "root:root1234@tcp(localhost:3306)/foo")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Queryx("SELECT * FROM Vedrfolnir LIMIT 10000")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	totalTime := time.Duration(0)

	i := 0
	for rows.Next() {
		var obs_id int
		var patient_id int
		var concept_name string
		var description string
		totalTimeAData := time.Duration(0)

		err := rows.Scan(&obs_id, &patient_id, &concept_name, &description)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("patient_id: %d, concept_name: %s, description: %s\n", patient_id, concept_name, description)

		//i, err := strconv.Atoi(patient_id)
		if HSMT[patient_id] == nil {
			//print("this is an empty user\n")
			HSMT[patient_id] = core.InitIndex(key0)
		}

		keyword := strings.Split(description, " ")
		conceptNameKeyword := strings.Split(concept_name, " ")
		keyword = append(keyword, conceptNameKeyword...)
		deDupKeyword := RemoveDuplicates(keyword)

		for _, kw := range deDupKeyword {

			kwHashStr := string(VedCrypto.Hash([]byte(kw)))
			lastChar := []byte{kwHashStr[len(kwHashStr)-1]}
			tempProof, _ := VedCrypto.Encrypt(lastChar, key0)
			nonexistenceProof := VedCrypto.Hash(tempProof, tempProof)

			proofs, res, _ := core.CheckExistence([]byte(kw), HSMT[patient_id], nonexistenceProof)
			if res == false {
				//print("keyword ", kw, " doest not exist, create an BHT_CMT \n")
				//time2_1 := time.Now()
				BHT_CMT := core.NewIndex(key1, 1)
				time2_2 := time.Now()
				HSMT[patient_id] = core.UpdateIndex(proofs, BHT_CMT)
				time2_3 := time.Now()
				totalTimeAData += time2_3.Sub(time2_2)
				DataAllPatient[patient_id] = append(DataAllPatient[patient_id], kw)
			}
		}

		totalTime += totalTimeAData
		i++
		if i%10000 == 0 {
			fmt.Printf("Average time cost for each data in %d data: %s \n", i, totalTime/time.Duration(i))
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func TestIndexInit() {
	key := key0
	time1 := time.Now()
	for i := 0; i < 10000; i++ {
		key[i%16] = key[0] >> i
		HSMT[i] = core.InitIndex(key)
	}

	time2 := time.Now()
	for i := 0; i < 10000; i++ {
		index := core.NewIndex(key1, 19)
		proofs, _ := core.GetProof([]byte("apple"), HSMT[i])
		HSMT[i] = core.UpdateIndex(proofs, index)
	}

	time3 := time.Now()

	initTime := time2.Sub(time1)
	updateTime := time3.Sub(time2)
	fmt.Printf("Mutiple index init %s, and update %s \n", initTime, updateTime)
}

func TestIndexVerify(patient_id int, kw []byte) {

	time1 := time.Now()
	//lastChar := []byte{kw[len(kw)-1]}
	// existProof, _ := VedCrypto.Encrypt(lastChar, key0)
	proof, leaf := core.GetProof(kw, HSMT[patient_id])
	//hash := VedCrypto.Hash(existProof, existProof)
	core.VerifyProof(proof, HSMT[0], leaf)

	time2 := time.Now()

	verifyTime := time2.Sub(time1)

	fmt.Printf("One index verify: %s \n", verifyTime)
}

func TestIndexVerifyTimeSpace() {
	db, err := sqlx.Open("mysql", "root:root1234@tcp(localhost:3306)/foo")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Queryx("SELECT * FROM Vedrfolnir LIMIT 6000000")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var obs_id int
		var patient_id int
		var concept_name string
		var description string
		// 将查询结果映射到结构体或变量
		err := rows.Scan(&obs_id, &patient_id, &concept_name, &description)
		if err != nil {
			log.Fatal(err)
		}

		// 处理每一行的数据
		// fmt.Printf("patient_id: %d, concept_name: %s, description: %s\n", patient_id, concept_name, description)

		keyword := strings.Split(description, " ")
		conceptNameKeyword := strings.Split(concept_name, " ")
		keyword = append(keyword, conceptNameKeyword...)
		deDupKeyword := RemoveDuplicates(keyword)
		Data[patient_id] += len(deDupKeyword)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	var maxPatient int
	var maxValue int

	for key, value := range Data {
		if maxValue < value {
			maxPatient = key
			maxValue = value
		}
	}

	print("patient id: ", maxPatient, " key: ", maxValue, "\n")
	HSMTaPatient := core.InitIndex(key0)

	rows, err = db.Queryx("SELECT * FROM Vedrfolnir LIMIT 6000000")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var obs_id int
		var patient_id int
		var concept_name string
		var description string
		err := rows.Scan(&obs_id, &patient_id, &concept_name, &description)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("patient_id: %d, concept_name: %s, description: %s\n", patient_id, concept_name, description)

		if patient_id == maxPatient {
			keyword := strings.Split(description, " ")
			conceptNameKeyword := strings.Split(concept_name, " ")
			keyword = append(keyword, conceptNameKeyword...)
			deDupKeyword := RemoveDuplicates(keyword)

			for _, kw := range deDupKeyword {

				kwHashStr := string(VedCrypto.Hash([]byte(kw)))
				lastChar := []byte{kwHashStr[len(kwHashStr)-1]}
				tempProof, _ := VedCrypto.Encrypt(lastChar, key0)
				nonexistenceProof := VedCrypto.Hash(tempProof, tempProof)

				proofs, res, _ := core.CheckExistence([]byte(kw), HSMTaPatient, nonexistenceProof)
				if res == false {
					//print("keyword ", kw, " doest not exist, create an BHT_CMT \n")
					//time2_1 := time.Now()
					BHT_CMT := core.NewIndex(key1, 1)
					HSMTaPatient = core.UpdateIndex(proofs, BHT_CMT)
					DataAPatient = append(DataAPatient, kw)
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	var proofStorage uintptr
	var keyStorage uintptr
	var verifyTime = time.Duration(0)
	for _, kw := range DataAPatient {
		kwHashStr := string(VedCrypto.Hash([]byte(kw)))
		lastChar := []byte{kwHashStr[len(kwHashStr)-1]}
		tempProof, _ := VedCrypto.Encrypt(lastChar, key0)
		nonexistenceProof := VedCrypto.Hash(tempProof, tempProof)

		proofs, _, target := core.CheckExistence([]byte(kw), HSMTaPatient, nonexistenceProof)
		time0 := time.Now()
		core.VerifyProof(proofs, HSMTaPatient, target)
		time1 := time.Now()
		verifyTime += time1.Sub(time0)
		proofStorage += unsafe.Sizeof(proofs)
		proofStorage += unsafe.Sizeof(target)
		keyStorage += unsafe.Sizeof(core.NewIndex(key1, 1))
	}
	fmt.Printf("Verify time:  %s, and storage cost: %d bytes. \n", verifyTime, proofStorage+keyStorage)
}

func TestIndexVerifyTimeSpaceAll() {
	var proofStorage uintptr
	var keyStorage uintptr
	var verifyTime = time.Duration(0)
	for patient_id, aHSMT := range HSMT {
		var proofStorageApatient uintptr
		var keyStorageApatient uintptr
		var verifyTimeApatient = time.Duration(0)
		for _, kw := range DataAllPatient[patient_id] {
			kwHashStr := string(VedCrypto.Hash([]byte(kw)))
			lastChar := []byte{kwHashStr[len(kwHashStr)-1]}
			tempProof, _ := VedCrypto.Encrypt(lastChar, key0)
			nonexistenceProof := VedCrypto.Hash(tempProof, tempProof)

			proofs, _, target := core.CheckExistence([]byte(kw), aHSMT, nonexistenceProof)
			time0 := time.Now()
			core.VerifyProof(proofs, aHSMT, target)
			time1 := time.Now()
			verifyTimeApatient += time1.Sub(time0)
			proofStorageApatient += unsafe.Sizeof(proofs)
			proofStorageApatient += unsafe.Sizeof(target)
			keyStorageApatient += unsafe.Sizeof(core.NewIndex(key1, 1))
		}
		//fmt.Printf("Patient: %d, verify time:  %s, and storage cost: %d bytes. \n", patient_id, verifyTimeApatient, proofStorageApatient+keyStorageApatient)
		proofStorage += proofStorageApatient
		keyStorage += keyStorageApatient
		verifyTime += verifyTimeApatient
	}
	fmt.Printf("Verify time:  %s, and storage cost: %d (proof) + %d (key) =  %d bytes. \n", verifyTime, proofStorage, keyStorage, proofStorage+keyStorage)

}

func RemoveDuplicates(slice []string) []string {
	encountered := map[string]bool{}
	var result []string

	for _, v := range slice {
		if encountered[v] == false {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
