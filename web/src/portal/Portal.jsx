import { React, useEffect, useState } from "react";
import { Grid, GridItem, Box, Center, Link, VStack, Spacer, Flex, Input } from "@chakra-ui/react"
import { Text } from "@chakra-ui/react"
import Server from '../server';
import Login from './Login';
import Desktop from '../desktop/Desktop';
import { GiGrapes } from "react-icons/gi";
import { GrNewWindow } from "react-icons/gr";
import { FiUsers } from "react-icons/fi";


function Portal() {

    const [token, setToken] = useState(localStorage.token);
    const [activeProject, setActiveProject] = useState(localStorage.project || null);
    const [projects, setProjects] = useState([]);

    function getProjectsList() {
        token && Server.getProjectsList()
            .then(setProjects)
    }
    useEffect(getProjectsList, [token]);

    function selectProject(project) {
        if (project) {
            localStorage.setItem("project", project)
        } else {
            localStorage.removeItem("project")
        }
        setActiveProject(project)
    }

    const projectsBoxes = projects.map(project => <Link onClick={_ => selectProject(project)}>
        <Box key={project} w="10em" h="10em" bg="blue.500" >
            <VStack textAlign="center" spacing="24px" >
                <Spacer />
                <GiGrapes size="50" />
                <Text color="white" isTruncated>{project}</Text>
            </VStack>

        </Box>
    </Link>)

    projectsBoxes.push(<Link >
        <Box w="10em" h="10em" bg="red.500" >
            <VStack textAlign="center" spacing="24px" >
                <Spacer />
                <GrNewWindow size="50" />
                <Text color="white" isTruncated>New Project</Text>
            </VStack>
        </Box>
    </Link>)

    projectsBoxes.push(<Link >
        <Box w="10em" h="10em" bg="yellow.500" >
            <VStack textAlign="center" spacing="24px" >
                <Spacer />
                <FiUsers size="50" />
                <Text color="white" isTruncated>Users</Text>
            </VStack>
        </Box>
    </Link>)


    const content = activeProject ? <Desktop project={activeProject} onExit={_ => selectProject(null)} /> : <Center>
        <Flex align="center"
            direction="column"            
            justify="space-between"
            wrap="wrap"
            
        >
            <Spacer minHeight="10px"/>
            <Input type="text" w="40%"></Input>
            <Spacer minHeight="200px"/>
            <Grid templateColumns="repeat(5, 1fr)" gap={3} >
                {projectsBoxes}
            </Grid>
            <Spacer />
        </Flex>
    </Center>


    return <>
        <Login isOpen={!token} onToken={setToken} />
        {content}
    </>



}
export default Portal