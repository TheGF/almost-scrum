import { Badge, Box } from "@chakra-ui/react";
import { React } from "react";
import MarkdownView from '../core/MarkdownView';


function TaskViewer(props) {
    const { task, searchKeys, saveTask, readOnly, height } = props;

    function onChange(value) {
      if (task) {
        task.description = value;
        saveTask(task);
      }
    }

    return task && <Box style={{ overflow: 'auto' }}>
        <div className="mde-preview">
            <div className="mde-preview-content task-viewer" style={{ height: height, overflow: 'auto' }}>
                <MarkdownView highlights={searchKeys}
                    readOnly={readOnly}
                    value={task.description}
                    onChange={onChange}
                     />
            </div>
        </div>
    </Box>

}

export default TaskViewer