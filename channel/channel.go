package channel

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
)

type ChannelLevel struct {
	ProbabilityError float64 // Вероятность ошибки
	ProbabilityLoss  float64 // Вероятность потери сообщения
}

func computeParityBit(encoded []int, positions []int) int {
	sum := 0
	for _, pos := range positions {
		sum += int(encoded[pos])
	}
	return int(sum % 2)
}

func EncodeHamming74(data []int) []int {

	encoded := make([]int, 7)
	encoded[0] = data[0]
	encoded[1] = data[1]
	encoded[2] = data[2]
	encoded[4] = data[3]

	encoded[3] = computeParityBit(encoded, []int{0, 1, 2})
	encoded[5] = computeParityBit(encoded, []int{0, 1, 4})
	encoded[6] = computeParityBit(encoded, []int{0, 2, 4})

	return encoded
}

func EncodeData(data []int) []int {

	encodedMessage := make([]int, 0)

	for i := 0; i < len(data); i += 4 {
		block := data[i:min(i+4, len(data))]   // Блок размером 4 байта
		encodedBlock := EncodeHamming74(block) // Кодируем блок
		encodedMessage = append(encodedMessage, encodedBlock...)
	}

	return encodedMessage
}

func setError(data []int) []int {
	errorPos := rand.Intn(len(data))
	data[errorPos] = (data[errorPos] + 1) % 2
	return data
}

func DecodeHamming74(data []int) []int {
	decodedData := make([]int, 4)
	syndrome := make([]int, 3)
	syndrome[0] = (data[0] + data[1] + data[2] + data[3]) % 2
	syndrome[1] = (data[0] + data[1] + data[4] + data[5]) % 2
	syndrome[2] = (data[0] + data[2] + data[4] + data[6]) % 2
	if (syndrome[0] + syndrome[1] + syndrome[2]) > 0 {
		pos := syndrome[0]*4 + syndrome[1]*2 + syndrome[2]
		data[int(math.Abs(float64(pos-7)))] = (data[int(math.Abs(float64(pos-7)))] + 1) % 2
	}
	decodedData[0] = data[0]
	decodedData[1] = data[1]
	decodedData[2] = data[2]
	decodedData[3] = data[4]

	return decodedData
}

func DecodeData(data []int) []int {

	decodedMessage := make([]int, 0)

	for i := 0; i < len(data); i += 7 {
		block := data[i:min(i+7, len(data))]
		decodedBlock := DecodeHamming74(block)
		decodedMessage = append(decodedMessage, decodedBlock...)
	}

	return decodedMessage
}

func bytesToBinary(bytes []byte) string {
	binaryStr := ""
	for _, b := range bytes {
		binaryStr += fmt.Sprintf("%08b", b)
	}
	return binaryStr
}

func bitsStringToBytes(bitsStr string) []byte {
	numBits := len(bitsStr)
	numBytes := (numBits + 7) / 8 // Округляем до ближайшего целого числа байтов
	bytes := make([]byte, numBytes)

	for i := 0; i < numBytes; i++ {
		for j := 0; j < 8; j++ {
			index := i*8 + j
			if index < numBits && bitsStr[index] == '1' {
				bytes[i] |= 1 << (7 - j) // Устанавливаем j-ый бит в i-том байте в 1
			}
		}
	}

	return bytes
}

func intArrayToString(intArray []int) string {
	var bitString string
	for _, num := range intArray {
		bitString += strconv.Itoa(num)
	}
	return bitString
}

func stringToIntArray(s string) []int {
	intArray := make([]int, len(s))
	for i, c := range s {
		intArray[i], _ = strconv.Atoi(string(c))
	}
	return intArray
}

func Channeltransmit(data []byte) ([]byte, bool, bool) {
	errorFlag := false
	channel := ChannelLevel{
		ProbabilityError: 11,
		ProbabilityLoss:  1,
	}
	binaryData := bytesToBinary(data)
	log.Println("binaryData", binaryData)
	bitArray := stringToIntArray(binaryData)
	encodedData := EncodeData(bitArray)

	if rand.Intn(100) < int(channel.ProbabilityError) {
		log.Println("errorFlag")
		encodedData = setError(encodedData)
		errorFlag = true
	}
	log.Println("decodedData")
	decodedData := DecodeData(encodedData)
	log.Println("finish")
	str := intArrayToString(decodedData)
	result := bitsStringToBytes(str)
	if rand.Intn(100) < int(channel.ProbabilityLoss) {
		return result, false, errorFlag
	}

	return result, true, errorFlag
}
