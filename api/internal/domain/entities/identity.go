package entities

import gonanoid "github.com/jaevor/go-nanoid"

func NewID() (string, error) {
   id, err := gonanoid.Standard(10)
   if err != nil {
       return "", err
   }
   return id(), nil
}