import { React, useEffect, useState } from "react";
import { Grid, GridItem, Box, Center, Link, VStack, Spacer, Flex, Input } from "@chakra-ui/react"
import { Text } from "@chakra-ui/react"
import Server from '../server';
import Login from './Login';
import Users from './Users';
import Desktop from '../desktop/Desktop';
import { GiGrapes } from "react-icons/gi";
import { GrNewWindow } from "react-icons/gr";
import { FiUsers } from "react-icons/fi";
import AddProject from './AddProject';


function Portal() {

    const [token, setToken] = useState(localStorage.token);
    const [activeProject, setActiveProject] = useState(null);
    const [projects, setProjects] = useState([]);
    const [showUsers, setShowUsers] = useState(false)
    const [showNewProject, setShowNewProject] = useState(false)

    function getProjectsList() {
        token && Server.getProjectsList()
            .then(projects => {
                setProjects(projects);
                if (projects.includes(localStorage.project)) {
                    setActiveProject(localStorage.project);
                }
            })
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

    function onCreate(project) {
        setProjects([project, ...projects])
        setShowNewProject(false)
    }

    const projectsBoxes = projects.map(project => <Link key={project}
        onClick={_ => selectProject(project)}>
        <Box w="12em" h="12em" bg="blue.500" >
            <VStack textAlign="center" spacing="24px" >
                <Spacer />
                <GiGrapes size="50" />
                <Text w="11em" color="white" isTruncated>{project}</Text>
            </VStack>

        </Box>
    </Link>)

    projectsBoxes.push(<Link key="#newProject"
        onClick={_ => setShowNewProject(true)}  >
        <Box w="12em" h="12em" bg="red.500" >
            <VStack textAlign="center" spacing="24px" >
                <Spacer />
                <GrNewWindow size="50" />
                <Text color="white" isTruncated>Add Project</Text>
            </VStack>
        </Box>
    </Link>)

    projectsBoxes.push(<Link key="#users"
        onClick={_ => setShowUsers(true)} >
        <Box w="12em" h="12em" bg="yellow.500" >
            <VStack textAlign="center" spacing="24px" >
                <Spacer />
                <FiUsers size="50" />
                <Text color="white" isTruncated>Users</Text>
            </VStack>
        </Box>
    </Link>)


    const content = token && activeProject ?
        <Desktop project={activeProject} onExit={_ => selectProject(null)} /> :
        <Center>
            <Flex align="center"
                direction="column"
                justify="space-between"
                wrap="wrap"

            >
                <Spacer minHeight="10px" />
                <Input type="text" w="40%"></Input>
                <Spacer minHeight="200px" />
                <Grid templateColumns="repeat(5, 1fr)" gap={3} >
                    {projectsBoxes}
                </Grid>
                <Spacer />
            </Flex>
        </Center>


    return <>
        <Login isOpen={!token} onToken={setToken} />
        <Users isOpen={showUsers} onClose={_ => setShowUsers(false)} />
        <AddProject isOpen={showNewProject} onCreate={onCreate} onClose={_ => setShowNewProject(false)} />
        {content}
    </>



}
export default Portal