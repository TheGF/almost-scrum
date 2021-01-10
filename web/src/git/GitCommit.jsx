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
    const [infos, setInfos] = useState([]);

    return <Center>
        <HStack spacing={5}>
            <Button size="lg" colorScheme="blue">Commit</Button>
            <Button size="lg" colorScheme="green">Push</Button>
        </HStack>
    </Center>
}
export default GitCommit;