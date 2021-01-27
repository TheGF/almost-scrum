import { Badge, Box } from "@chakra-ui/react";
import { React } from "react";
import MarkdownView from '../core/MarkdownView';


function TaskViewer(props) {
    const { task, searchKeys, height } = props;

    return task && <Box style={{ overflow: 'auto' }}>
        <div className="mde-preview">
            <div className="mde-preview-content task-viewer" style={{ height: height, overflow: 'auto' }}>
                <MarkdownView highlights={searchKeys}
                    source={task.description} />
            </div>
        </div>
    </Box>

}

export default TaskViewer