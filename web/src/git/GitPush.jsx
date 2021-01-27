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
    const [pushOutput, setPushOutput] = useState(null)

    function push() {
        setPushInProgress(true)
        Server.postGitPush(project)
            .then(setPushOutput)
            .then(_ => setPushInProgress(false))
    }

    return pushOutput ?
        <VStack>
            <Text fontSize="md" color="green">
                Push was successful
            </Text>
            <Textarea
                value={pushOutput}
                size="md"
                resize="Vertical"
                rows="10"
            />
        </VStack> :
        <VStack>
            <Button size="lg" colorScheme="blue" isLoading={pushInProgress}
                onClick={push}>
                Push
            </Button>
        </VStack>

}
export default GitPush;