package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

func watchLog(file string, out io.Writer) {
	var err error
	var f *os.File

	f, err = os.Open(file)
	eFatal(err)
	defer f.Close()

	r := bufio.NewReader(f)
	info, err := f.Stat()
	eFatal(err)

	oldSize := info.Size()

	db, err := dbConnect()
	eFatal(err)
	defer db.Close()

	// goto is not harmful if it's used in a clear logic like this
	// This code block runs continuously so one more for {} would only increase scope making harder to understand all the logic
mainloop:
	err = db.Ping()
	eFatal(err)

	// ReadLine loop
	var errRL error
	var line []byte
	for {
		// Channel loop
		insChan := make(chan dbInsertResult, 10)
		for i := 0; i < 10; i++ {
			line, _, errRL = r.ReadLine()

			if errRL != nil {
				break
			}

			var cl clf
			cl, err = parse(string(line))
			if err != nil {
				log.Println(err)
				log.Println(" > ", string(line))
				continue
			}
			if cl.Method != "GET" && cl.Method != "POST" {
				log.Println("Parsed line is neither GET nor POST")
				continue
			}

			wg.Add(1)
			go dbInsertLog(insChan, cl, db)
		}

		if errRL == io.EOF {
			break
		}

		wg.Wait()
		close(insChan)

		// Return channel loop
		for res := range insChan {
			if res.err != nil {
				log.Println(res.err)
				log.Printf("> %+v", res.log)
			}
		}
	}

	pos, err := f.Seek(0, io.SeekCurrent)
	eFatal(err)

	// sleep and seek loop
	for {
		t := time.Minute * 5
		log.Printf("Watch sleeping for %v", t)
		time.Sleep(t)
		log.Println("Watch waking up")
		newinfo, err := f.Stat()
		eFatal(err)

		newSize := newinfo.Size()
		if newSize != oldSize {
			if newSize < oldSize {
				_, err = f.Seek(0, 0)
				eFatal(err)
			} else {
				_, err = f.Seek(pos, io.SeekStart)
				eFatal(err)
			}
			r = bufio.NewReader(f)
			oldSize = newSize
			break
		}
	}

	goto mainloop
}
