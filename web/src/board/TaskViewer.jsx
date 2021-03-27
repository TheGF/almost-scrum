import { Badge, Box, HStack } from "@chakra-ui/react";
import { React } from "react";
import MarkdownView from '../core/MarkdownView';
import Properties from './Properties';


function TaskViewer(props) {
  const { task, searchKeys, saveTask, readOnly, height, startEdit } = props;

  function onChange(value) {
    if (task) {
      task.description = value;
      saveTask(task);
    }
  }

  return task && <Box style={{ overflow: 'auto' }} onClick={startEdit}>
    <HStack>
      <Box w="90%" h={height}>
        <div className="mde-preview">
          <div className="mde-preview-content task-viewer" style={{
            height: height-41, overflow: 'auto'
          }}>
            <MarkdownView highlights={searchKeys}
              readOnly={readOnly}
              value={task.description}
              onChange={onChange}
            />
          </div>
        </div>
      </Box>
      <Box minW="240px" h={height - 20}>
        <Properties task={task} saveTask={saveTask} users={[]} height={height}
          height={height} readOnly={true} />
      </Box>
    </HStack>
  </Box>

}

export default TaskViewer