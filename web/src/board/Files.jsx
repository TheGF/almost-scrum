import {
  Center, Input, Select, Switch, Table, TableCaption, Tbody,
  Td, Tr
} from '@chakra-ui/react';
import { Tabs, TabList, TabPanels, Tab, TabPanel } from "@chakra-ui/react"
import { React, useContext, useState, useEffect } from "react";
import T from '../core/T';
import UserContext from '../UserContext';
import AttachedFiles from './AttachedFiles';
import Library from '../library/Library';

function Files(props) {
  const { task, saveTask, readOnly } = props;
  const { info } = useContext(UserContext)
  const [attachedFiles, setAttachedFiles] = useState(task.files || [])
  const [tabIndex, setTabIndex] = useState(0)
  const [path, setPath] = useState(0)

  function showPath(path) {
    setPath(path)
    setTabIndex(1)
  }

  function saveFiles() {
    if (task.files != attachedFiles) {
      saveTask({
        ...task,
        files: attachedFiles,
      })
    }
  }
  useEffect(saveFiles, [attachedFiles])

  return <Tabs variant="soft-rounded" size="sm" index={tabIndex} onChange={setTabIndex}>
    <TabList>
      <Tab>Attached</Tab>
      <Tab>Full Library</Tab>
    </TabList>

    <TabPanels>
      <TabPanel>
        <AttachedFiles attachedFiles={attachedFiles} setAttachedFiles={setAttachedFiles}
          onShowPath={showPath} readOnly={readOnly} />
      </TabPanel>
      <TabPanel>
        <Library attachedFiles={attachedFiles} setAttachedFiles={!readOnly && setAttachedFiles}
          path={path} />
      </TabPanel>
    </TabPanels>
  </Tabs>
}
export default Files;