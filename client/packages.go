package client

import (
   "io/ioutil"
   //"log"
   "net/http"
)

type RDepotConfig struct {
  Host string
  Token string
}

func PackagesList(cfg RDepotConfig) ([]byte, error) {

  req, err := http.NewRequest(
    "GET",
    cfg.Host + "/api/manager/packages/list",
    nil)

  if err != nil { return nil, err }

  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Bearer " + cfg.Token)

  resp, err := http.DefaultClient.Do(req)

  switch {
  case err != nil:
    return nil, err
  case resp.StatusCode != 200:
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
  default:
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
  }

}

