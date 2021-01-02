import { React, useEffect, useState, useContext } from "react";
import ReactMde from "react-mde";
import ReactMarkdown from "react-markdown";
import "react-mde/lib/styles/css/react-mde-all.css";
import { Box, Button } from "@chakra-ui/react";
import T from "../core/T";


function loadSuggestions(text) {
  return new Promise((accept, reject) => {
    setTimeout(() => {
      const suggestions = [
        {
          preview: "Andre",
          value: "@andre"
        },
        {
          preview: "Angela",
          value: "@angela"
        },
        {
          preview: "David",
          value: "@david"
        },
        {
          preview: "Louise",
          value: "@louise"
        }
      ].filter((i) => i.preview.toLowerCase().includes(text.toLowerCase()));
      accept(suggestions);
    }, 250);
  });
}



function TaskEditor(props) {
  const { name } = props;

  const save = async function* (data) {
    // Promise that waits for "time" milliseconds
    const wait = function (time) {
      return new Promise((a, r) => {
        setTimeout(() => a(), time);
      });
    };

    // Upload "data" to your server
    // Use XMLHttpRequest.send to send a FormData object containing
    // "data"
    // Check this question: https://stackoverflow.com/questions/18055422/how-to-receive-php-image-data-over-copy-n-paste-javascript-with-xmlhttprequest

    await wait(2000);
    // yields the URL that should be inserted in the markdown
    yield "https://picsum.photos/300";
    await wait(2000);

    // returns true meaning that the save was successful
    return true;
  };

  const { task, saveTask } = props
  const [value, setValue] = useState(task && task.description);

  const features = {
    name: "features",
    icon: () => (
      <Button colorScheme="yellow"><T>features</T></Button>
    ),
    execute: opts => {
      opts.textApi.replaceSelection("NICE");
    }
  };
  
  function onChange(value) {
    setValue(value);
    if (task) {
      task.description = value;
      saveTask(task);
    }
  }

  return (<Box>
    <ReactMde
      value={value}
      onChange={onChange}
      disablePreview={true}
      loadSuggestions={loadSuggestions}
      paste={{
        saveImage: save
      }}
    /></Box>
  );
}

export default TaskEditor