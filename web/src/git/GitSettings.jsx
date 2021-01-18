import {
    Box,
    Button,
    Input,
    InputGroup,
    InputRightElement,
    VStack,
    Center
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
import { React, useContext, useState } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function GitSettings(props) {
    const { project, info } = useContext(UserContext)
    const [showPassword, setShowPassword] = useState(false)

    return <Accordion>
        <AccordionItem>
            <AccordionButton>
                <Box flex="1" textAlign="left">
                    Password Based Authentication
                </Box>
                <AccordionIcon />
            </AccordionButton>
            <AccordionPanel pb={4}>
                <Center>
                <FormControl id="credentials" maxWidth="30em">
                    <FormLabel>Username</FormLabel>
                    <Input type="username" />
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
                </Center>
            </AccordionPanel>
        </AccordionItem>

        <AccordionItem>
            <AccordionButton>
                <Box flex="1" textAlign="left">
                    Section 2 title
        </Box>
                <AccordionIcon />
            </AccordionButton>
            <AccordionPanel pb={4}>


            </AccordionPanel>
        </AccordionItem>
    </Accordion>
}
export default GitSettings;