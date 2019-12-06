package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//There codes were written by Noah Chang

//OutputFile3DinCSV takes in a 3D slice of Matrix, []Matrix, and outputs a csvfile with a name matching the index of the input.
//It outputs in the folder "outputcsv3D" under the current directory, removing all the contents under the folder before writing the file.
func OutputFile3DinCSV(timepoints []Matrix) {

	folderName := "outputcsv3D"

	//Getting directory
	outputFolder := GetNewFolderDir(folderName)

	//If the directory does not exist, make one
	MakeDirIfNotExist(outputFolder)

	//Delete the previous csv files
	RefreshDirectory(outputFolder)

	for i := range timepoints {

		//Get filename for ith generation
		filename := outputFolder + "/3D_Matrix_" + strconv.Itoa(i) + ".csv"
		csvfile, err := os.Create(filename)
		if err != nil {
			fmt.Println("Couldn’t create the file!")
		}
		defer csvfile.Close()

		//Setting col names
		output := [][]string{{"x", "y", "z", "state"}}

		for x := range timepoints[i] {
			for y := range timepoints[i][x] {
				for z := range timepoints[i][x][y] {

					if timepoints[i][x][y][z].state != "h" {
						outputCoordinate := make([]string, 0)

						outputCoordinate = append(outputCoordinate, strconv.Itoa(x))
						outputCoordinate = append(outputCoordinate, strconv.Itoa(y))
						outputCoordinate = append(outputCoordinate, strconv.Itoa(z))

						//replacing the states with appropriate hex color codes
						if timepoints[i][x][y][z].state == "C" {
							outputCoordinate = append(outputCoordinate, "#ADD8E6")
						}

						//Options for further color coding----------------------------------
						// if timepoints[i][x][y][z].state == "Q" {
						// 	outputCoordinate = append(outputCoordinate, "#FFFF00")
						// }
						// if timepoints[i][x][y][z].state == "N" {
						// 	outputCoordinate = append(outputCoordinate, "#8B0000")
						// }
						// if timepoints[i][x][y][z].state == "wN" {
						// 	outputCoordinate = append(outputCoordinate, "#696969")
						// }
						//------------------------------------------------------------------
						output = append(output, outputCoordinate)
					}
				}
			}
		}

		//writing csv files
		writer := csv.NewWriter(csvfile)
		for _, elements := range output {
			err := writer.Write(elements)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		writer.Flush()
	}
}

//OutputFile2DinCSV functions similar to OutputFile3DinCSV.
//It takes in a 2D slice of Matrix, []Matrix, and outputs a csvfile with a name matching the index of the input.
//It outputs in the folder "outputcsv2D" under the current directory, removing all the contents under the folder before writing the file.
func OutputFile2DinCSV(timepoints []Matrix2D) {
	folderName := "outputcsv2D"

	outputFolder := GetNewFolderDir(folderName)

	MakeDirIfNotExist(outputFolder)

	RefreshDirectory(outputFolder)

	for i := range timepoints {

		filename := outputFolder + "/2D_Matrix_" + strconv.Itoa(i) + ".csv"
		csvfile, err := os.Create(filename)
		if err != nil {
			fmt.Println("Couldn’t create the file!")
		}
		defer csvfile.Close()

		output := [][]string{{"x", "y", "state"}}

		for x := range timepoints[i] {
			for y := range timepoints[i][x] {

				if timepoints[i][x][y].state != "h" {
					outputCoordinate := make([]string, 0)
					outputCoordinate = append(outputCoordinate, strconv.Itoa(x))
					outputCoordinate = append(outputCoordinate, strconv.Itoa(y))
					outputCoordinate = append(outputCoordinate, timepoints[i][x][y].state)
					output = append(output, outputCoordinate)
				}

			}
		}

		writer := csv.NewWriter(csvfile)
		for _, elements := range output {
			err := writer.Write(elements)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
		writer.Flush()
	}
}

//RefreshDirectory takes in a directory string and removes all *.csv content under it.
//Edited code from https://stackoverflow.com/questions/33450980/how-to-remove-all-contents-of-a-directory-using-golang
func RefreshDirectory(dir string) {
	d, _ := ioutil.ReadDir(dir)
	for _, files := range d {
		if strings.Contains(files.Name(), ".csv") {
			os.Remove(path.Join([]string{dir, files.Name()}...))
		}
	}
}

//MakeDirIfNotExist takes in a directory string and makes the directory(folder) if the directory does not exist.
//Edited code from https://siongui.github.io/2017/03/28/go-create-directory-if-not-exist/
func MakeDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic("Problem when making a new folder")
		}
	}
}

//GetNewFolderDir takes in a folderName string and creates a folder using the input string under the current directory.
func GetNewFolderDir(folderName string) string {
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	outputFolder := currentDirectory + "/" + folderName

	return outputFolder
}

//ReadPNGs takes in a directory and reads all the ".png" files.
//It does not goes into subfolders
//It returns []image.Image which can later be made as a gif
func ReadPNGs(dir string) []image.Image {

	//reads the files under the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	files = SortByFileName(files)

	//Declaring the list of images
	imageList := make([]image.Image, 0)

	for _, file := range files {
		//Setting the file directory
		fileDir := dir + "/" + file.Name()

		//only reads in if the file has .png extension
		if filepath.Ext(fileDir) == ".png" {
			pngFile, err := os.Open(fileDir)

			if err != nil {
				panic("Error while reading in png file")
			}

			defer pngFile.Close()

			//Decoding the png file
			src, _, err := image.Decode(pngFile)
			if err != nil {
				panic("Error while decoding image")
			}

			//Appending to the list of images
			imageList = append(imageList, src)
		}
	}

	return imageList
}

//SortByFileName takes in a list of files and sort them.
func SortByFileName(files []os.FileInfo) []os.FileInfo {

	PNGOnly := make([]os.FileInfo, 0)

	for _, i := range files {
		length := len(i.Name())
		if i.Name()[length-4:] == ".png" {
			PNGOnly = append(PNGOnly, i)
		}
	}
	fmt.Println(len(PNGOnly))

	filenumbers := make([]string, 0)

	for i := range PNGOnly {
		filenumbers = append(filenumbers, strconv.Itoa(i))
	}

	sortedFiles := make([]os.FileInfo, 0)

	for i := range filenumbers {
		for _, j := range PNGOnly {

			//indexing is specific! if any change in naming in R should correspond here too!
			if filenumbers[i] == j.Name()[13:len(j.Name())-5] {

				sortedFiles = append(sortedFiles, j)

			}
		}
	}

	for _, i := range sortedFiles {
		fmt.Println(i.Name())
	}
	return sortedFiles
}

//OutputFileMetastasisInCSV writes CSV according to the slices of gens of metastasis cell counts
func OutputFileMetastasisInCSV(metaSlice [][3]int) {

	filename := "metastasis.csv"
	csvfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Couldn’t create the file!")
	}
	defer csvfile.Close()

	//Naming the columns
	output := [][]string{{"Bones", "Lungs", "Liver"}}

	for i := range metaSlice {
		string := make([]string, 0)

		//Append to each columns
		string = append(string, strconv.Itoa(metaSlice[i][0]))
		string = append(string, strconv.Itoa(metaSlice[i][1]))
		string = append(string, strconv.Itoa(metaSlice[i][2]))

		output = append(output, string)
	}

	//Writing csv files...
	writer := csv.NewWriter(csvfile)
	for _, elements := range output {
		err := writer.Write(elements)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	writer.Flush()

}
