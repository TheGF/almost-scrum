import {
    Box,
    Button,
    Input,
    InputGroup,
    InputRightElement,
    VStack,
    Center,
    StackDivider,
    Spacer
} from "@chakra-ui/react";
import { Tabs, TabList, TabPanels, Tab, TabPanel } from "@chakra-ui/react"
import {
    Accordion,
    AccordionItem,
    AccordionButton,
    AccordionPanel,
    AccordionIcon,
} from "@chakra-ui/react"
import {
    FormControl,
    FormLabel,
    FormErrorMessage,
    FormHelperText,
} from "@chakra-ui/react"
import { Switch } from "@chakra-ui/react"
import { React, useContext, useState, useEffect } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function GitSettings(props) {
    const { project, info } = useContext(UserContext)
    const [showPassword, setShowPassword] = useState(false)
    const [gitSettings, setGitSettings] = useState(null)

    function getGitSettings() {
        Server.getGitSettings(project)
            .then(setGitSettings)
    }
    useEffect(getGitSettings, [])

    function switchNativeGit() {
        setGitSettings({
            ...gitSettings,
            useGitNative: !gitSettings.useGitNative,
        })
        Server.putGitSettings(gitSettings)
    }

    return gitSettings ? <VStack spacing={10}>
        <FormControl id="git-native-form" maxWidth="30em">
            <FormLabel>Use Git Native Client</FormLabel>
            <Switch id="git-native" isChecked={gitSettings.useGitNative} onChange={switchNativeGit} />
            <FormHelperText>
                Requires Git client to be installed in your system. When off, connect directly to Git (not recommended)
        </FormHelperText>
        </FormControl>
        <FormControl id="credentials" maxWidth="30em">
            <FormLabel>Username</FormLabel>
            <Input type="username" value={gitSettings.username}/>
            <FormLabel>Password</FormLabel>
            <InputGroup size="md">
                <Input
                    pr="4.5rem"
                    type={showPassword ? "text" : "password"}
                    placeholder="Enter password"
                />
                <InputRightElement width="4.5rem">
                    <Button h="1.75rem" size="sm" onClick={
                        _ => setShowPassword(!showPassword)
                    }>
                        {showPassword ? "Hide" : "Show"}
                    </Button>
                </InputRightElement>
            </InputGroup>
            <FormHelperText>We'll encrypt your password</FormHelperText>
        </FormControl>
    </VStack > : null
}
export default GitSettings;