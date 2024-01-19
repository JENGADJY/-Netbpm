package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Définition de la structure PBM pour représenter une image PBM
type PBM struct {
	data          [][]bool // Les données binaires (true pour 1, false pour 0)
	width, height int      // Largeur et hauteur de l'image
	magicNumber   string   // Numéro magique pour identifier le type de fichier PBM
}

// Fonction pour lire un fichier PBM et créer une instance PBM
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Lecture du numéro magique PBM
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, errors.New("Numéro magique PBM invalide")
	}

	// Lecture des dimensions de l'image
	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, errors.New("Dimensions invalides")
	}

	// Conversion des dimensions en entiers
	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, errors.New("Largeur invalide")
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, errors.New("Hauteur invalide")
	}

	// Initialisation des données PBM
	data := make([][]bool, height)
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		if magicNumber == "P1" {
			data[i] = parseP1Line(line, width)
		} else {
			data[i] = parseP4Line(line, width)
		}
	}

	// Création et renvoi de l'objet PBM
	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Fonction pour analyser une ligne P1 et créer un tableau booléen correspondant
func parseP1Line(line string, width int) []bool {
	data := make([]bool, width)

	// Séparation de la ligne en caractères individuels
	chars := strings.Fields(line)

	// Vérification si le nombre de caractères correspond à la largeur attendue
	if len(chars) != width {
		return nil
	}

	// Conversion des caractères en booléens
	for i, char := range chars {
		data[i] = char == "1"
	}
	return data
}

// Fonction pour analyser une ligne P4 et créer un tableau booléen correspondant
func parseP4Line(line string, width int) []bool {
	data := make([]bool, width)

	// Vérification que la ligne a suffisamment d'octets pour couvrir la largeur
	if len(line) < (width+7)/8 {
		return nil
	}

	// Parcours de la largeur de l'image
	for i := 0; i < width; i++ {
		// Calcul de l'index d'octet et de la position du bit dans l'octet
		byteIndex := i / 8
		bitPos := uint(7 - (i % 8))

		// Vérification que byteIndex est dans les limites de la ligne
		if byteIndex >= len(line) {
			return nil
		}

		// Extraction du bit de l'octet
		bit := (line[byteIndex] >> bitPos) & 1
		data[i] = bit == 1
	}

	return data
}

// Méthode pour obtenir la taille de l'image PBM
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// Méthode pour obtenir la valeur d'un pixel à une position spécifique dans l'image PBM
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Méthode pour définir la valeur d'un pixel à une position spécifique dans l'image PBM
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Méthode pour sauvegarder l'image PBM dans un fichier
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écriture du numéro magique et des dimensions de l'image dans le fichier
	_, err = fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)
	if err != nil {
		return err
	}

	// Boucle pour écrire les données binaires dans le fichier
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pbm.magicNumber == "P1" {
				if pixel {
					_, err = file.WriteString("1 ")
				} else {
					_, err = file.WriteString("0 ")
				}
			} else {
				if pixel {
					_, err = file.Write([]byte{0xFF})
				} else {
					_, err = file.Write([]byte{0x00})
				}
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Méthode pour inverser les couleurs de l'image PBM
func (pbm *PBM) Invert() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Méthode pour inverser les lignes de l'image PBM
func (pbm *PBM) Flip() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width/2; x++ {
			pbm.data[y][x], pbm.data[y][pbm.width-x-1] = pbm.data[y][pbm.width-x-1], pbm.data[y][x]
		}
	}
}

// Méthode pour inverser les colonnes de l'image PBM
func (pbm *PBM) Flop() {
	for y := 0; y < pbm.height/2; y++ {
		pbm.data[y], pbm.data[pbm.height-y-1] = pbm.data[pbm.height-y-1], pbm.data[y]
	}
}

// Méthode pour ajouter le première elementdu fichier
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func main() {
	filename := "C:/Users/JENGO/Netbpm/PBM/duck.pbm" // mon chemin
	pbm, err := ReadPBM(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Magic Number: %s\n", pbm.magicNumber)
	fmt.Printf("Width: %d\n", pbm.width)
	fmt.Printf("Height: %d\n", pbm.height)

	width, height := pbm.Size()
	fmt.Printf("Image Size: %d x %d\n", width, height)

	x, y := 2, 3
	fmt.Printf("Pixel at (%d, %d): %v\n", x, y, pbm.At(x, y))

	newValue := true
	pbm.Set(x, y, newValue)
	fmt.Printf("New pixel value at (%d, %d): %v\n", x, y, pbm.At(x, y))

	outputFilename := "output.pbm"
	err = pbm.Save(outputFilename)
	if err != nil {
		fmt.Println("Error saving the PBM image:", err)
		return
	}

	fmt.Println("PBM image saved successfully to", outputFilename)
}
