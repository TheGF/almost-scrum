import {
    Button, Flex, HStack, Switch, Table, Tbody, Td, Th, Thead, Tr,
    VStack
} from '@chakra-ui/react';
import { React, useContext } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function GitFiles(props) {
    const { project } = useContext(UserContext)
    const { gitStatus, setGitStatus } = props;

    function fetchStatus() {
        Server.getGitStatus(project)
            .then(setGitStatus)
    }

    function getRows(gitStatus) {
        if (!gitStatus) {
            return [null, null];
        }
        const { stagedFiles, untrackedFiles } = gitStatus;

        function switchFileStage(file) {
            const idx = stagedFiles.indexOf(file);
            if (idx != -1) {
                stagedFiles.splice(idx, 1);
                untrackedFiles.push(file);
            } else {
                const idx = untrackedFiles.indexOf(file);
                untrackedFiles.splice(idx, 1);
                stagedFiles.push(file);
            }
            setGitStatus({
                ...gitStatus,
                stagedFiles: stagedFiles,
                untrackedFiles: untrackedFiles
            });
        }

        const staged = stagedFiles && stagedFiles.sort().map(file => <Tr>
            <Td>{file}</Td>
            <Td><Switch isChecked={true}
                onChange={_ => switchFileStage(file)} /></Td>
        </Tr>);

        const untracked = untrackedFiles && untrackedFiles.sort().map(file => <Tr>
            <Td>{file}</Td>
            <Td><Switch isChecked={false}
                onChange={_ => switchFileStage(file)} /></Td>
        </Tr>);

        return [staged, untracked];
    }

    const [staged, untracked] = getRows(gitStatus)

    return <VStack>
        <HStack>
            <Button onClick={fetchStatus}>Get Status</Button>
            <Button>Pull</Button>
        </HStack>
        <Flex overflow="auto" h="20em" w="70%">
            <Table overflow="auto" >
                <Thead>
                    <Tr>
                        <Th>File</Th>
                        <Th>Staged</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    {staged}
                    {untracked}
                </Tbody>
            </Table>
        </Flex>
    </VStack>
}
export default GitFiles;