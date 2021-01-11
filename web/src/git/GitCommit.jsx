import {
    Accordion,
    AccordionButton,
    AccordionIcon, AccordionItem,
    AccordionPanel, Box, Button, Center, Spacer,
    HStack, List, ListIcon, ListItem, StackDivider, Textarea, VStack, Flex, Input, Text,
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import { BiCheckCircle, BiCircle } from "react-icons/bi";
import Server from '../server';
import UserContext from '../UserContext';

function GitCommit(props) {
    const { project, info } = useContext(UserContext)
    const { gitStatus, gitMessage } = props
    const [gitHash, setGitHash] = useState(null)
    const commitInfo = {
        user: info.loginUser,
        header: gitMessage.header,
        body: gitMessage.body,
        files: gitStatus && [...gitStatus.ashFiles, ...gitStatus.stagedFiles] || [],
    }

    function commit() {
        Server.postGitCommit(project, commitInfo)
            .then(setGitHash)
    }

    return <VStack>
        <VStack textAlign="left">
            <b>Summary</b>
            <label>User: {commitInfo.user}</label>
            <label>Header: {commitInfo.header}</label>
            <label>Staged Files: {commitInfo.files.join('')}</label>
        </VStack>
        <HStack spacing={5}>
            <Button size="lg" colorScheme="blue" onClick={commit}>Commit</Button>
            <Button size="lg" colorScheme="green">Push</Button>
        </HStack>
    </VStack>

}
export default GitCommit;