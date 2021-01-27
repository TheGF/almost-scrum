import {
    Button,
    Textarea,
    VStack,
    Text
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import Server from '../server';
import UserContext from '../UserContext';

function GitPush(props) {
    const { project, info } = useContext(UserContext)
    const [pushInProgress, setPushInProgress] = useState(false)
    const [pushOutput, setPushOutput] = useState(null)

    function push() {
        setPushInProgress(true)
        Server.postGitPush(project)
            .then(setPushOutput)
            .then(_ => setPushInProgress(false))
    }

    return <VStack>
        <Button size="lg" colorScheme="blue" isLoading={pushInProgress}
            onClick={push}>
            Push
        </Button>
        {pushOutput ? <>
            <Textarea
                value={pushOutput}
                size="md"
                resize="Vertical"
                rows="6"
            />
        </> : null}
    </VStack>

}
export default GitPush;