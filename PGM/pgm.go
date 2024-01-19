package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Définition de la structure PGM pour représenter une image PGM
type PGM struct {
	data        [][]uint8 // Données de l'image (valeurs de pixels)
	width       int       // Largeur de l'image
	height      int       // Hauteur de l'image
	magicNumber string    // Numéro magique pour identifier le type de fichier PGM
	max         int       // Valeur maximale autorisée pour un pixel
}

// Fonction pour lire un fichier PGM et créer une instance PGM
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, errors.New("Format PGM non pris en charge")
	}

	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(dimensions[0])
	height, _ := strconv.Atoi(dimensions[1])

	scanner.Scan()
	maxVal, _ := strconv.Atoi(scanner.Text())

	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
		for j := range data[i] {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			data[i][j] = uint8(val)
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         maxVal,
	}, nil
}

// Méthode pour obtenir la taille de l'image PGM
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// Méthode pour obtenir la valeur d'un pixel à une position spécifique dans l'image PGM
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Méthode pour définir la valeur d'un pixel à une position spécifique dans l'image PGM
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Méthode pour sauvegarder l'image PGM dans un fichier
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	if pgm.magicNumber == "P2" {
		for i := 0; i < pgm.height; i++ {
			for j := 0; j < pgm.width; j++ {
				fmt.Fprintf(writer, "%d ", pgm.data[i][j])
			}
			fmt.Fprintln(writer)
		}
	} else if pgm.magicNumber == "P5" {
		for i := 0; i < pgm.height; i++ {
			for j := 0; j < pgm.width; j++ {
				fmt.Fprintf(writer, "%c", pgm.data[i][j])
			}
		}
	} else {
		return errors.New("Format PGM non pris en charge")
	}

	return writer.Flush()
}

// Méthode pour inverser les couleurs de l'image PGM
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Méthode pour inverser les lignes de l'image PGM
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Méthode pour inverser les colonnes de l'image PGM
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
	}
}

// Méthode pour définir le numéro magique de l'image PGM
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// Méthode pour définir la valeur maximale autorisée pour un pixel dans l'image PGM
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Méthode pour faire pivoter l'image PGM de 90 degrés dans le sens des aiguilles d'une montre
func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for i := range newData {
		newData[i] = make([]uint8, pgm.height)
		for j := range newData[i] {
			newData[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// Définition de la structure PBM pour représenter une image PBM
type PBM struct {
	data          [][]bool // Données binaires de l'image (true pour 1, false pour 0)
	width, height int      // Largeur et hauteur de l'image
	magicNumber   string   // Numéro magique pour identifier le type de fichier PBM
}

// Méthode pour convertir une image PGM en une image PBM
func (pgm *PGM) ToPBM() *PBM {
	pbmData := make([][]bool, pgm.height)
	for i := 0; i < pgm.height; i++ {
		pbmData[i] = make([]bool, pgm.width)
		for j := 0; j < pgm.width; j++ {
			pbmData[i][j] = pgm.data[i][j] > uint8(pgm.max)/2
		}
	}
	return &PBM{
		data:        pbmData,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}
}

func main() {
	//mon chemin
	pgmFilename := "C:/Users/JENGO/Netbpm/PGM/duck.pgm"
	pgm, err := ReadPGM(pgmFilename)
	if err != nil {
		fmt.Println("Error reading PGM:", err)
		return
	}

	fmt.Printf("PGM Magic Number: %s\n", pgm.magicNumber)
	fmt.Printf("PGM Width: %d\n", pgm.width)
	fmt.Printf("PGM Height: %d\n", pgm.height)
	fmt.Printf("PGM Max Value: %d\n", pgm.max)

	pgm.Invert()

	pgm.Flip()

	//mon chemin
	modifiedPGMFilename := "C:/Users/JENGO/Netbpm/PGM/duck_edit.pgm"
	err = pgm.Save(modifiedPGMFilename)
	if err != nil {
		fmt.Println("Error saving modified PGM:", err)
		return
	}
	fmt.Println("Modified PGM image saved successfully to", modifiedPGMFilename)

	pbm := pgm.ToPBM()

	fmt.Printf("\nPBM Magic Number: %s\n", pbm.magicNumber)
	fmt.Printf("PBM Width: %d\n", pbm.width)
	fmt.Printf("PBM Height: %d\n", pbm.height)

}
