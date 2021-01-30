import {
    Button, Flex, HStack, Switch, Table, Tbody, Td, Th, Thead, Tr,
    VStack
} from '@chakra-ui/react';
import { React, useContext, useState } from "react";
import T from '../core/T';
import Server from '../server';
import UserContext from '../UserContext';

const change2descr = {
    'A': 'added',
    'M': 'modified',
    'R': 'renamed',
    'D': 'deleted',
    '?': 'untracked',
}

function GitFiles(props) {
    const { project } = useContext(UserContext)
    const { stagedFiles, setStagedFiles } = props
    const [gitStatus, setGitStatus] = useState(null)
    const [gettingStatus, setGettingStatus] = useState(false);
    const [pullInProgress, setPullInProgress] = useState(false);

    function getStatus() {
        setGettingStatus(true)
        Server.getGitStatus(project)
            .then(status => {
                const stagedFiles = status ? Object.keys(status.files)
                    .filter(file => status.files[file] != '?')
                    .concat(status.ashFiles) : []

                setGitStatus(status)
                setStagedFiles(stagedFiles)
                setGettingStatus(false)
            })
    }

    function pull() {
        setPullInProgress(true)
        Server.postGitPull(project)
            .then(_ => setPullInProgress(false))
    }

    function getRows(gitStatus) {
        if (!gitStatus) {
            return [];
        }

        function switchFileStage(file) {
            const idx = stagedFiles.indexOf(file);
            if (idx != -1) {
                stagedFiles.splice(idx, 1);
            } else {
                stagedFiles.push(file);
            }
            setStagedFiles(stagedFiles);
        }

        return Object.keys(gitStatus.files).sort().map(file => <Tr>
            <Td>{file}</Td>
            <Td><T>{change2descr[gitStatus.files[file]]}</T></Td>
            <Td><Switch isChecked={stagedFiles.includes(file)}
                onChange={_ => switchFileStage(file)} /></Td>
        </Tr>);
    }

    const rows = getRows(gitStatus)

    return <VStack>
        <HStack>
            <Button onClick={getStatus} isLoading={gettingStatus}>Get Status</Button>
            <Button onClick={pull} isLoading={pullInProgress}>Pull</Button>
        </HStack>
        <Flex overflow="auto" h="20em" w="100%">
            <Table overflow="auto" >
                <Thead>
                    <Tr>
                        <Th>File</Th>
                        <Th>Change</Th>
                        <Th>Staged</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {rows}
                </Tbody>
            </Table>
        </Flex>
    </VStack>
}
export default GitFiles;