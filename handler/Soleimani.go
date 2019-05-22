package handler

import (
	"bufio"
	"fmt"
	"github.com/lupengyu/trafficflow/constant"
	"github.com/lupengyu/trafficflow/helper"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func culLAG(positions []constant.PositionMeta) {
	L := 0.0
	A := 0.0
	G0 := 0.0 //Latitude
	G1 := 0.0 //Longitude
	for i := 1; i < len(positions); i++ {
		position1 := &constant.Position{
			Longitude: positions[i-1].Longitude,
			Latitude:  positions[i-1].Latitude,
		}
		position2 := &constant.Position{
			Longitude: positions[i].Longitude,
			Latitude:  positions[i].Latitude,
		}
		L += helper.PositionSpacing(position1, position2)
		A += math.Abs(position2.Longitude-position1.Longitude) * (position2.Latitude + position1.Latitude) * 0.5
		G0 += math.Abs(position2.Latitude - position1.Latitude)
		G1 += math.Abs(position2.Longitude - position1.Longitude)
	}
	fmt.Println("L=", L)
	fmt.Println("A=", A)
	fmt.Println("G0=", G0)
	fmt.Println("G1=", G1)
}

func Soleimani(reFile string, opFile string) {
	file, err := os.Open("data/" + reFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	reData := make([]constant.PositionMeta, 0)
	bfRd := bufio.NewReader(file)
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		strs := strings.Split(string(line), ",")
		if len(strs) < 4 {
			continue
		}
		longitude, _ := strconv.ParseFloat(strs[0], 64)
		latitude, _ := strconv.ParseFloat(strs[1], 64)
		reData = append(reData, constant.PositionMeta{
			Longitude: longitude,
			Latitude:  latitude,
		})
	}

	file, err = os.Open("data/" + opFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()
	opData := make([]constant.PositionMeta, 0)
	bfRd = bufio.NewReader(file)
	for {
		line, _, _ := bfRd.ReadLine()
		if len(line) == 0 {
			break
		}
		strs := strings.Split(string(line), ",")
		if len(strs) < 2 {
			continue
		}
		longitude, _ := strconv.ParseFloat(strs[0], 64)
		latitude, _ := strconv.ParseFloat(strs[1], 64)
		opData = append(opData, constant.PositionMeta{
			Longitude: longitude,
			Latitude:  latitude,
		})
	}

	fmt.Println("re=================")
	culLAG(reData)
	fmt.Println("op=================")
	culLAG(opData)

	fmt.Println("Score=================")
	Score1 := (21445.727247748928 - 5880.780155638942) / 21445.727247748928
	Score2 := 111319.0 / 21445.727247748928 * (1.7637562436624175 - 0.6354144695830733 + 0.1784470012090189 - 0.045065166968800696 + 0.07208910942942737 - 0.02597105008290157)
	fmt.Println(Score1, Score2, Score1+Score2)
}
