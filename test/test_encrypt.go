package main

import (
	core "Ved/core"
	"Ved/lib/VedCrypto"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
	"time"
)

func main_temp() {
	//TestSECorrectness()
	//TestEncryptOptimizationCorrectness()
	TestEncrypt()
	TestSearch()
}
func TestSECorrectness() {
	keyword1 := "this is my first key word"
	keyword2 := "this is my first key word"
	timestamp := "2022.1"
	key := core.CreateKey([]byte(keyword1), []byte(timestamp))
	print(len(key), "\n")
	ciphertext := core.Encryption([]byte(keyword1), key)
	res := core.Search(ciphertext, []byte(keyword2), key)
	print(res)
}

func TestEncryptOptimization() {
	keyword1 := "apple"
	keyword2 := "apple"
	timestamp := 2
	//keyRaw := core.CreateKey_One([]byte(keyword1), timestamp)
	keyRaw := core.CreateKey_All([]byte(keyword1))
	key := keyRaw[timestamp][:32]
	//print(len(key), "\n")
	ciphertext := core.Encryption([]byte(keyword1), key)
	res := core.Search(ciphertext, []byte(keyword2), key)
	print(res)
}

func TestEncrypt() {

	db, err := sqlx.Open("mysql", "root:root1234@tcp(localhost:3306)/foo")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Queryx("SELECT * FROM Vedrfolnir")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	stmt, err := db.Prepare("INSERT INTO patient_data (patient_id, ciphertext, hash) VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	i := 0
	totalTime := time.Duration(0)
	for rows.Next() {
		var obs_id int
		var patient_id int
		var description string
		one_record_encrypt_time := time.Duration(0)

		err := rows.Scan(&obs_id, &patient_id, &description)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("patient_id: %d, concept_name: %s, description: %s\n", patient_id, concept_name, description)

		//i, err := strconv.Atoi(patient_id)

		keyword := strings.Split(description, " ")
		deDupKeyword := removeDuplicates(keyword)
		hash := VedCrypto.Hash([]byte(description))

		/*
			var ciphertext []byte
			for _, kw := range deDupKeyword {
				time0 := time.Now()
				one_word_ciphertext := core.Encryption([]byte(kw), []byte("ThisIs16ByteKey1"))
				time1 := time.Now()
				//print(kw, " ", one_word_ciphertext, "\n")
				//print(len(one_word_ciphertext), "\n")
				ciphertext = append(ciphertext, one_word_ciphertext...)
				one_record_encrypt_time += time1.Sub(time0)
			}
			_, err = stmt.Exec(patient_id, ciphertext, hash)
			if err != nil {
				panic(err.Error())
			}
			totalTime += one_record_encrypt_time

		*/
		//var ciphertext []byte

		for _, kw := range deDupKeyword {
			time0 := time.Now()
			one_word_ciphertext := core.Encryption([]byte(kw), []byte("ThisIs16ByteKey1"))
			time1 := time.Now()
			//print(kw, " ", one_word_ciphertext, "\n")
			//print(len(one_word_ciphertext), "\n")
			_, err = stmt.Exec(patient_id, one_word_ciphertext, hash)
			if err != nil {
				panic(err.Error())
			}
			one_record_encrypt_time += time1.Sub(time0)
		}

		totalTime += one_record_encrypt_time

		i++
		if i%10000 == 0 {
			fmt.Printf("Average time cost for each data in %d data: %s \n", i, totalTime/time.Duration(i))
		}
	}
	fmt.Printf("Total encryption time cost for %d data: %s \n", i, totalTime)
}

func TestSearch() {
	db, err := sqlx.Open("mysql", "root:root1234@tcp(localhost:3306)/foo")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Queryx("SELECT * FROM patient_data") //select count(*) from patient_data;
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	totalTime := time.Duration(0)
	i := 0
	for rows.Next() {
		//var patient_id int
		var patient_id int
		var ciphertextRaw []byte
		var hash []byte
		err := rows.Scan(&patient_id, &ciphertextRaw, &hash)
		if err != nil {
			log.Fatal(err)
		}
		ciphertext := splitBytes(ciphertextRaw)

		key := []byte("ThisIs16ByteKey1")
		for _, ct := range ciphertext {
			kw := []byte("VISIT")
			time0 := time.Now()
			core.Search(ct, kw, key)
			time1 := time.Now()

			totalTime += time1.Sub(time0)
			//print(" ", res)
		}
		i++
		if i%10000 == 0 {
			fmt.Printf("Average time cost for each data in %d data: %s \n", i, totalTime/time.Duration(i))
		}
	}
	fmt.Printf("Total search time cost for %d data: %s \n", i, totalTime)
}

func removeDuplicates(slice []string) []string {
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

func splitBytes(input []byte) [][]byte {
	var result [][]byte

	for i := 0; i < len(input); i += 64 {
		end := i + 64
		if end > len(input) {
			end = len(input)
		}
		result = append(result, input[i:end])
	}

	return result
}
