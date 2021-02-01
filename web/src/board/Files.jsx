import {
  Center, Input, Select, Switch, Table, TableCaption, Tbody,
  Td, Tr, Flex
} from '@chakra-ui/react';
import { Tabs, TabList, TabPanels, Tab, TabPanel } from "@chakra-ui/react"
import { React, useContext, useState, useEffect } from "react";
import T from '../core/T';
import UserContext from '../UserContext';
import AttachedFiles from './AttachedFiles';
import Library from '../library/Library';

function Files(props) {
  const { task, saveTask, readOnly, height } = props;
  const [attachedFiles, setAttachedFiles] = useState(task.files || [])
  const [tabIndex, setTabIndex] = useState(0)
  const [path, setPath] = useState(0)

  function showPath(path) {
    setPath(path)
    setTabIndex(1)
  }

  function setAttachedFilesAndSave(files) {
    setAttachedFiles(files)
    saveTask({
      ...task,
      files: files,
    })
  }

  return <Tabs variant="soft-rounded" size="sm" index={tabIndex}
      onChange={setTabIndex} isLazy>
      <TabList>
        <Tab>Attached</Tab>
        <Tab>Full Library</Tab>
      </TabList>

      <TabPanels>
        <TabPanel>
          <AttachedFiles attachedFiles={attachedFiles} setAttachedFiles={setAttachedFilesAndSave}
            onShowPath={showPath} readOnly={readOnly} height={height}/>
        </TabPanel>
        <TabPanel>
          <Library attachedFiles={attachedFiles} height={height}
            setAttachedFiles={!readOnly && setAttachedFilesAndSave}
            path={path} />
        </TabPanel>
      </TabPanels>
    </Tabs>
}
export default Files;