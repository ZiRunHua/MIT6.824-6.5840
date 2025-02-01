package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Lab struct {
	Name, Short, Path string
	Tests             []string
}

var (
	RootDir = getRootDir()
	Labs    = []Lab{
		Lab3, Lab3A, Lab3B, Lab3C, Lab3D, Lab4, Lab4A, Lab4B,
	}
	Lab3 = Lab{
		"lab3",
		"3",
		filepath.Join("src", "raft"),
		append(Lab3A.Tests, append(Lab3B.Tests, append(Lab3C.Tests, Lab3D.Tests...)...)...),
	}
	Lab3A = Lab{
		"lab3A",
		"3A",
		filepath.Join("src", "raft"),
		[]string{
			"TestInitialElection3A",
			"TestReElection3A",
			"TestManyElections3A",
		},
	}
	Lab3B = Lab{
		"lab3B",
		"3B",
		filepath.Join("src", "raft"),
		[]string{
			"TestBasicAgree3B",
			"TestRPCBytes3B",
			"TestFollowerFailure3B",
			"TestLeaderFailure3B",
			"TestFailAgree3B",
			"TestFailNoAgree3B",
			"TestConcurrentStarts3B",
			"TestRejoin3B",
			"TestBackup3B",
			"TestCount3B",
		},
	}
	Lab3C = Lab{
		"lab3C",
		"3C",
		filepath.Join("src", "raft"),
		[]string{
			"TestPersist13C",
			"TestPersist23C",
			"TestPersist33C",
			"TestFigure83C",
			"TestUnreliableAgree3C",
			"TestFigure8Unreliable3C",
			"TestReliableChurn3C",
			"TestUnreliableChurn3C",
		},
	}
	Lab3D = Lab{
		"lab3D",
		"3D",
		filepath.Join("src", "raft"),
		[]string{
			"TestSnapshotBasic3D",
			"TestSnapshotInstall3D",
			"TestSnapshotInstallUnreliable3D",
			"TestSnapshotInstallCrash3D",
			"TestSnapshotInstallUnCrash3D",
			"TestSnapshotAllCrash3D",
			"TestSnapshotInit3D",
		},
	}
	Lab4 = Lab{
		"lab4",
		"4",
		filepath.Join("src", "kvraft"),
		append(Lab4A.Tests, Lab4B.Tests...),
	}
	Lab4A = Lab{
		"lab4A",
		"4A",
		filepath.Join("src", "kvraft"),
		[]string{
			"TestBasic4A",
			"TestSpeed4A",
			"TestConcurrent4A",
			"TestUnreliable4A",
			"TestUnreliableOneKey4A",
			"TestOnePartition4A",
			"TestManyPartitionsOneClient4A",
			"TestManyPartitionsManyClients4A",
			"TestPersistOneClient4A",
			"TestPersistConcurrent4A",
			"TestPersistConcurrentUnreliable4A",
			"TestPersistPartition4A",
			"TestPersistPartitionUnreliable4A",
			"TestPersistPartitionUnreliableLinearizable4A",
		},
	}
	Lab4B = Lab{
		"lab4B",
		"4B",
		filepath.Join("src", "kvraft"),
		[]string{
			"TestSnapshotRPC4B",
			"TestSnapshotSize4B",
			"TestSpeed4B",
			"TestSnapshotRecover4B",
			"TestSnapshotRecoverManyClients4B",
			"TestSnapshotUnreliable4B",
			"TestSnapshotUnreliableRecover4B",
			"TestSnapshotUnreliableRecoverConcurrentPartition4B",
			"TestSnapshotUnreliableRecoverConcurrentPartitionLinearizable4B",
		},
	}
	LabMap = makeLabMap()
)

func makeLabMap() map[string]*Lab {
	lm := make(map[string]*Lab, len(Labs))
	for i, lab := range Labs {
		lm[lab.Name] = &Labs[i]
	}
	return lm
}

func getRootDir() string {
	// `os.Getwd()` is avoided here because, during tests, the working directory is set to the test file’s directory.
	// This command retrieves the module's root directory instead.
	rootDir, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()
	if err == nil {
		return filepath.Dir(strings.TrimSpace(string(rootDir)))
	}
	// If `go list` fails, it may indicate the absence of a Go environment.
	// In such cases, this suggests we are not in a test environment, so fall back to `os.Getwd()` to set `RootDir`.
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Validate that the directory exists
	_, err = os.Stat(workDir)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("Path:%s does not exists", workDir))
		}
		panic(err)
	}
	// 当前目录是test
	return filepath.Dir(filepath.Clean(workDir))
}
