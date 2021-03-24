import {
    Box, Center, Flex, Grid,
    HStack, IconButton, Input, Link, Spacer, Text, VStack
} from "@chakra-ui/react";
import { React, useEffect, useState } from "react";
import { BiLockOpen } from "react-icons/bi";
import { FiHelpCircle } from "react-icons/fi";
import { GiGrapes } from "react-icons/gi";
import { GrLogout, GrNewWindow } from "react-icons/gr";
import Desktop from '../desktop/Desktop';
import Help from '../help/Help';
import Server from '../server';
import Access from './Access';
import AddProject from './AddProject';
import ChangePassword from './ChangePassword';
import ClaimInvite from './ClaimInvite';
import Login from './Login';


function Portal(props) {

    const { systemUser } = props
    const [token, setToken] = useState(localStorage.token);
    const [activeProject, setActiveProject] = useState(null);
    const [projects, setProjects] = useState([]);
    const [showUsers, setShowUsers] = useState(false)
    const [showNewProject, setShowNewProject] = useState(false)
    const [showHelp, setShowHelp] = useState(false);
    const [showChangePassword, setShowChangePassword] = useState(false)

    const [filter, setFilter] = useState('')

    function getProjectsList() {
        token && Server.getProjectsList()
            .then(projects => {
                if (projects && projects.includes) {
                    setProjects(projects);
                    if (projects.includes(localStorage.project)) {
                        setActiveProject(localStorage.project);
                    }
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

    function logout() {
        localStorage.removeItem("token")
        setToken(null)
    }

    function authenticated(token, weakPassword) {
        setToken(token)
        setShowChangePassword(weakPassword)
    }

    const projectsBoxes = projects ? projects
        .filter(p => p.toLowerCase().includes(filter.toLowerCase()))
        .map(project => <Link key={project}
            onClick={_ => selectProject(project)}>
            <Box w="12em" h="12em" bg="blue.500" >
                <VStack textAlign="center" spacing="24px" >
                    <Spacer />
                    <GiGrapes size="50" />
                    <Text w="11em" color="white" isTruncated>{project}</Text>
                </VStack>
            </Box>
        </Link>) : []

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
    
    projectsBoxes.push(<ClaimInvite key="claimInvite"/>)

    projectsBoxes.push(<Link key="#users"
        onClick={_ => setShowUsers(true)} >
        <Box w="12em" h="12em" bg="yellow.500" >
            <VStack textAlign="center" spacing="24px" >
                <Spacer />
                <BiLockOpen size="50" />
                <Text color="white" isTruncated>Access</Text>
            </VStack>
        </Box>
    </Link>)


    function getContent() {
        if (!token) return null;

        if (showChangePassword) {
            return <ChangePassword isOpen={showChangePassword}
                onClose={_ => setShowChangePassword(false)} />
        }

        if (activeProject && projects.includes(activeProject)) {
            return <Desktop key={activeProject} project={activeProject} onExit={_=>selectProject(null)} />
        }

        return <Center>
            <ChangePassword isOpen={showChangePassword}
                onClose={_ => setShowChangePassword(false)} />
            <Access isOpen={showUsers} onClose={_ => setShowUsers(false)} />
            <AddProject isOpen={showNewProject} onCreate={onCreate} onClose={_ => setShowNewProject(false)} />
            <Help isOpen={showHelp} onClose={_ => setShowHelp(false)} />
            <Flex align="center"
                direction="column"
                justify="space-between"
                wrap="wrap">
                <Spacer minHeight="10px" />
                <HStack w="55%">
                    <Input type="text" value={filter}
                        onChange={e => setFilter(e.target.value)} />
                    <IconButton icon={<FiHelpCircle />} title="Help" onClick={_ => setShowHelp(true)} />
                    <IconButton icon={<GrLogout />} title="Logout" onClick={logout} />
                </HStack>
                <Spacer minHeight="200px" />
                <Grid templateColumns="repeat(5, 1fr)" gap={3} >
                    {projectsBoxes}
                </Grid>
                <Spacer />
            </Flex>
        </Center>

    }

    document.title = activeProject ? `Almost Scrum: ${activeProject}` : 'Almost Scrum';
    return <>
        <Login isOpen={!token} systemUser={systemUser} onAuthenticated={authenticated} />
        {getContent()}
    </>



}
export default Portal