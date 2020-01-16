//Package cmd ...
/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"gansible/pkg/utils"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var commands string
var hosts string
var wg sync.WaitGroup
var timeout int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run commands on multiple hosts in parallel",
	Long: `Run commands on multiple hosts in parallel,return result when finished.Default number of concurrenrt tasks is 5.
Default timeout of each task is 300 seconds.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip, err := utils.ParseIPStr(hosts)
		if err != nil {
			fmt.Println(err)
		}
		if ip == nil {
			fmt.Println("No hosts specified!")
		} else {
			if forks < 1 {
				forks = 1
			} else if forks > 10000 {
				fmt.Println("Max forks is 10000")
				return
			}
			p, _ := ants.NewPool(forks)
			defer p.Release()
			for _, host := range ip {
				wg.Add(1)
				//_ = p.Submit(func() {
				//runr := utils.DoCommand(host, commands, timeout)
				//runinfo := utils.RunInfo(runr)
				//fmt.Println(runinfo)
				//wg.Done()
				//})
				_ = p.Submit(func() {
					passwords := []string{"abc", "passw0rd"}
					var client *ssh.Client
					client, _ = utils.TryPasswords("root", passwords, host, 22, 30)
					if client == nil {
						fmt.Println("All passwords are wrong.")
					} else {
						defer client.Close()
						timeout := 300
						execr := utils.Execute(client, commands, timeout)
						execinfo := utils.ExecInfo(host, execr)
						fmt.Println(execinfo)
						wg.Done()
					}
				})
			}
			wg.Wait()
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().StringVarP(&commands, "commands", "c", "", "separate multiple command with semicolons(eg: pwd;ls)")
	runCmd.Flags().StringVarP(&hosts, "hosts", "H", "", "eg: 10.0.0.1;10.0.0.2-5;10.0.0.6-10.0.0.8")
	runCmd.Flags().IntVarP(&timeout, "timeout", "", 300, "task should finished before timeout")
	runCmd.MarkFlagRequired("commands")
	runCmd.MarkFlagRequired("hosts")
}
