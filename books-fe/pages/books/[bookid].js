import axios from 'axios'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'

function UpdateBook({book}){
  
  const [id , setId] = useState(book.id)
  const [title , setTitle] = useState(book.title)
  const [isbnNo , setISBNo] = useState(book.ISBNno)
  const [series , setSeries] = useState(book.Series)

  function handleAddBook(e){
    e.preventDefault()
    axios.patch(`http://localhost:3000/api/books/${book.id}`, {
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


const EditBook = () => {
  const router = useRouter()
  const { bookid } = router.query

  if (!bookid){
    console.log("loading")
    return "Loading ID"
  }

  const [book, setBook] = useState(undefined)

  useEffect(()=>{
  
    axios.get(`/api/books/${bookid}`)
    .then(response=>
      {
        setBook(response.data)
        console.log(response)
      })
   .catch(e=>console.log("err: ",e))
  
  }, [bookid])

  

  if(!book){ 
    return "loading"
  }
  return <UpdateBook book={book}/>
}

export default EditBook