package resticcmd

import (
	"encoding/json"
	"fmt"
	"qrestic/types"
	"strings"
)

func executeCmd(cmdLine string, env ...string) ([]byte, error) {
	cmd := prepareCmd(cmdLine, env...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		fmt.Println(string(out))
		msg := fmt.Sprint(err)
		if len(out) > 0 {
			msg = string(out)
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return out, nil
}

type snapshot struct {
	Time     string
	Tree     string
	Paths    []string
	Hostname string
	Username string
	Uid      int
	Gid      int
	Id       string
	Short_id string
}

// used to ensure unique string content required by tree with string
func spaces(len int) string {
	out := ""
	ind := 0
	for ind < len {
		out += " "
		ind++
	}
	return out
}

func convertTime(time string) string {
	time = time[:strings.Index(time, ".")]
	return strings.ReplaceAll(time, "T", " ")
}

func GetSnapshots() (types.SnapshotTree, error) {
	if !configLoaded {
		panic("call ReadConfig first")
	}
	cmdLine := fmt.Sprintf("restic|-r|%s|snapshots|--json", configs[selectedConfig].Repository)
	data, err := executeCmd(cmdLine, "RESTIC_PASSWORD="+configs[selectedConfig].Password)
	if err != nil {
		return nil, err
	}
	var snaps []snapshot
	err = json.Unmarshal(data, &snaps)
	if err != nil {
		return nil, err
	}
	treeData := make(types.SnapshotTree)
	for ind, snap := range snaps {
		time := convertTime(snap.Time)
		treeData[""] = append(treeData[""], time)
		var leafData []string
		leafData = append(leafData, "Id: "+snap.Short_id+spaces(ind))
		leafData = append(leafData, "Host: "+snap.Hostname+spaces(ind))
		leafData = append(leafData, "User: "+snap.Username+spaces(ind))
		for _, path := range snap.Paths {
			leafData = append(leafData, "Path: "+path+spaces(ind))
		}
		treeData[time] = leafData
	}
	fmt.Println("done")
	return treeData, nil
}
