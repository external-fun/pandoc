import React, {useEffect, useState} from 'react';
import './App.css';
import { DragAndDrop } from './DragAndDrop';
import * as api from './api/api'

interface BasePage {
  type: 'UploadPage' | 'UploadingPage';
}

interface UploadPage extends BasePage {
  type: 'UploadPage';
}

interface UploadingPage extends BasePage {
  type: 'UploadingPage';
  fileName: string;
}

type Page = UploadPage | UploadingPage;

function App() {
  const [page, setPage] = useState<Page>({
    type: 'UploadPage'
  });
  const [uuid, setUuid] = useState<string>("")

  const uploadFile = async (file: File) => {
    const resp = await api.uploadFile(file, 'markdown', 'pdf')
    if (resp.type === "UploadResponse") {
      setUuid(resp.uuid)
    }
  };

  const addFile = async (file: File) => {
    setPage({
      type: 'UploadingPage',
      fileName: file.name
    });
    await uploadFile(file);
  };

  return (
    <div className="App centered">
      {page.type === 'UploadPage' && <DragAndDrop onAddFile={addFile} /> }
      {page.type === 'UploadingPage' && <span>{page.fileName}</span> }
    </div>
  );
}

export default App;
