import {
    Button,
    VStack
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function GitPush(props) {
    const { project, info } = useContext(UserContext)
    const [pushInProgress, setPushInProgress] = useState(false)

    function push() {
        setPushInProgress(true)
        Server.postGitPush(project)
            .then(_=>setPushInProgress(false))
    }

    return <VStack>
            <Button size="lg" colorScheme="blue" isLoading={pushInProgress}
                onClick={push}>
                Push
            </Button>
    </VStack>

}
export default GitPush;