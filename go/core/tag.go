package core


func ListTags(p Project) ([]string, error) {
	path := filepath.Join(project.Path, ProjectTagsFolder)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Warnf("Error reading path store")
		return nil, err
	}

	[]string tags = make([]string, len(fileInfos))
	for _, fileInfo := range fileInfos {
		tags = append(tags, fileInfo.Name)
	}
	return tags, nil
}

func ResolveTag