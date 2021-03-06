import React, { useEffect, useState } from 'react';
import { Modal } from 'reactstrap';

import axios from 'axios';
import { render } from 'react-dom';


function AddBook(props){
  
  const [id , setId] = useState()
  const [title , setTitle] = useState()
  const [isbnNo , setISBNo] = useState()
  const [series , setSeries] = useState()

  function handleAddBook(e){
    e.preventDefault()
    axios.post(`http://localhost:3000/api/books`, {
      "id" : parseInt(id), 
      "title": title, 
      "ISBNno":  isbnNo, 
      "Author":  "", 
      "Series": series, 
      "Genre": "", 
      "Rating": 3.8,
      
    }).then(r=>console.log("msg: ", r))
    .catch(e=>console.log("err: ",e))
    console.log("add", id, title, isbnNo, series)
    
  }
  


  return(
    <form >
    <label for="id">ID</label>
    <input  type="int" value={id} onChange={e=>setId(e.target.value)} name="id"/>
    <label for="title">TITLE</label>
    <input id="title" type="text" value={title} onChange={e=>setTitle(e.target.value)} name="title"></input>
    <label for="isbnno">ISBNNo</label>
    <input id="isbnno" type="text" value={isbnNo} onChange={e=>setISBNo(e.target.value)} name="isbno"></input>
    <label for="series">SERIES</label>
    <input id="series" type="text" value={series} onChange={e=>setSeries(e.target.value)} name="series"></input>
    <button onClick={handleAddBook}>AddBook</button>
   </form>
  )
}





function Book(props) {
  const book = props.book

  function handleDelete() {
    console.log(book)

    axios.delete(`http://localhost:3000/api/books/${book.id}`)
      .then((response) => {
        console.log("success")

      })
      .catch(error => {
        console.log(error);
        // this.setState({...this.state, isFetching: false});
      })

  }
  return (
    <tr key={book.id}>
      <td>{book.id}</td>
      <td>{book.title}</td>
      <td>{book.ISBNno}</td>
      <td>{book.Author}</td>
      <td>{book.Series}</td>
      <td>{book.Genre}</td>
      <td>

        <button  size="sm">UPDATE</button>
        <button onClick={handleDelete} size="sm">DELETE</button>

      </td>
    </tr>

  )

}

function BookList() {
  const [launchedBooks, setLaunchedBooks] = useState([]);
  
  // state = [newbookdata= {
  //   id: "",
  //   title: "",
  // }] 


  useEffect(() => {
    axios.get('http://localhost:3000/api/books')
      .then((response) => {
        setLaunchedBooks(response.data)
      })

      .catch(error => {
        console.log(error);
        // this.setState({...this.state, isFetching: false});
      })

  }, [])

  let books = launchedBooks.map((book) => (<Book book={book} />));

  return (
    <div className="App container">
    
     
      <table>
        <thead>
          <tr>
            <th>Id</th>
            <th>Title</th>
            <th>ISBNNo</th>
            <th>Author</th>
            <th>Series</th>
            <th>Genre</th>
            <th>Actions</th>
          </tr>
        </thead>

        <tbody>
          {books}
        </tbody>

      </table>
     
     );
    </div>
  );

}

export default function Home(){
  return(
    <>
    <AddBook/>
    <BookList/>
    </>
  )
}