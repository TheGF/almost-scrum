import {
    Button, Center, FormControl, FormHelperText, FormLabel, HStack, Input,
    InputGroup, InputRightElement, Spacer, Switch, VStack
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import T from "../core/T";
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
        const settings = {
            ...gitSettings,
            useGitNative: !gitSettings.useGitNative,
        }
        setGitSettings(settings)
    }

    function setUsername(e) {
        const username = e && e.target.value
        const settings = {
            ...gitSettings,
            username: username,
        }
        setGitSettings(settings)
    }

    function setPassword(e) {
        const password = e && e.target.value
        const settings = {
            ...gitSettings,
            password: password,
        }
        setGitSettings(settings)
    }

    function saveSettings() {
        Server.putGitSettings(project, settings)
    }

    return gitSettings ? <VStack spacing={10}>
        <FormControl id="git-native-form" maxWidth="30em">
            <HStack>
                <FormLabel>Use Git Native Client</FormLabel>
                <Spacer />
                <Switch id="git-native" isChecked={gitSettings.useGitNative} onChange={switchNativeGit} />
            </HStack>
            <FormHelperText>
                Requires Git client to be installed in your system. When off, connect directly to Git (not recommended)
        </FormHelperText>
        </FormControl>
        <FormControl id="credentials" isRequired maxWidth="30em">
            <FormLabel>Username</FormLabel>
            <Input type="username"
                value={gitSettings.username}
                onChange={setUsername} />
            <FormLabel>Password</FormLabel>
            <InputGroup size="md">
                <Input
                    pr="4.5rem"
                    value={gitSettings.password}
                    type={showPassword ? "text" : "password"}
                    onChange={setPassword}
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
            <FormHelperText><T>We'll encrypt your password</T></FormHelperText>
        </FormControl>
        <Center>
            <Button onClick={saveSettings}><T>save settings</T></Button>
        </Center>
    </VStack > : null
}
export default GitSettings;