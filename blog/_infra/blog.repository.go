package blog_infra

import (
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/db"
)

type BlogMemoryRepository struct {
	base db.BaseMemoryRepository[blog_models.Blog]
}

func (repo *BlogMemoryRepository) Create(blog blog_models.Blog) error {
	return repo.base.Create(blog)
}

func (repo *BlogMemoryRepository) List() ([]blog_models.Blog, error) {
	return repo.base.List()
}

func (repo *BlogMemoryRepository) Get(id string) (blog_models.Blog, error) {
	return repo.base.Get(id)
}

func (repo *BlogMemoryRepository) Update(blog blog_models.Blog) error {
	return repo.base.Update(blog)
}

func (repo *BlogMemoryRepository) Delete(id string) error {
	return repo.base.Delete(id)
}

func CreateMemoryBlogRepo() *BlogMemoryRepository {
	repo := &BlogMemoryRepository{base: *db.CreateMemoryRepo[blog_models.Blog]()}
	repo.seed()
	return repo
}

var snippet = "```"

func (repo *BlogMemoryRepository) seed() {
	entities := []blog_models.Blog{blog_models.BlogFromJSON([]byte(`{"Id":"523cfd06-6283-4d3a-937a-5dc68ffcb9ed","Title":"Understanding coupling and cohesion","Body":"In the realm of software engineering, two fundamental concepts that stand as pillars for designing robust, maintainable, and scalable systems are coupling and cohesion. Coupling refers to the degree of direct knowledge or reliance one class or module has on another. This means that changes in one module could directly affect another module if the coupling is tight, leading to a fragile system architecture where changes can cascade and cause unforeseen issues. Conversely, low coupling or loose coupling aims for modules that are more independent, making the system more modular and easier to modify or extend. \n\n Cohesion, on the other hand, deals with how closely related and focused the responsibilities of a single module are. High cohesion within a module means that its functions are closely related to each other and perform a well-defined task, making the system more understandable and easier to maintain. It encourages the breaking down of a system into manageable pieces, where each piece is responsible for a single, focused task or group of related tasks. High cohesion and low coupling are often seen as ideals in software design, leading to systems that are easier to debug, test, and maintain over time. \n\n Understanding and applying the principles of coupling and cohesion can significantly impact the quality of software systems. By striving for low coupling, developers ensure that modules can be developed, tested, and debugged in isolation, reducing complexity and dependency issues. Similarly, by aiming for high cohesion within modules, developers can create more intuitive and manageable codebases. These principles guide the design of software that is not only functional but also adaptable and resilient to change, qualities that are essential in the fast-paced and ever-evolving field of software development.","RelatedPosts":[{"Id":"49e6a869-6713-43d1-9431-45f491fea19f","Title":"The Principles of SOLID: Building Blocks for Object-Oriented Design"},{"Id":"82a40beb-bb7e-4874-8bd3-ffac983e555c","Title":"Refactoring Techniques for Reducing Coupling in Your Code"},{"Id":"767bdd8b-9972-4912-beb1-948cfcad4ae2","Title":"Design Patterns: Strategies for Enhancing Cohesion"}]}`)),
		blog_models.BlogFromJSON([]byte(`{"Id":"767bdd8b-9972-4912-beb1-948cfcad4ae2","Title":"Design Patterns: Strategies for Enhancing Cohesion","Body":"In the dynamic world of software development, **Design Patterns** play a pivotal role in crafting efficient and maintainable code. 🚀 They are not just solutions to common problems but strategies that enhance the cohesion of your codebase, making it more modular and easier to understand. Cohesion refers to how closely related and focused the responsibilities of a single module are. High cohesion within modules or components means that they have a well-defined purpose and are less likely to be affected by changes in other parts of the system. 🎯 This principle is at the heart of many design patterns, which aim to encapsulate changes, promote single responsibility, and ultimately, create a more robust and adaptable codebase. 🧩\n\nOne classic example is the **Strategy Pattern**. This pattern enables an object to change its behavior at runtime by making algorithms interchangeable within that object's family. By leveraging the Strategy Pattern, developers can increase cohesion by separating the concerns of how an object performs a certain task from the object itself. 🔄 This separation of concerns not only enhances code maintainability but also promotes a more declarative style of programming, where the focus is on the 'what' rather than the 'how.'\n\nAnother powerful tool in the arsenal of design patterns is the **Decorator Pattern**. It allows for behavior to be added to individual objects, either statically or dynamically, without affecting the behavior of other objects from the same class. This is particularly useful in adhering to the Open-Closed Principle, which states that software entities should be open for extension, but closed for modification. 🛠️ With the Decorator Pattern, cohesion is enhanced by enabling a more flexible extension of object functionality without necessitating changes to existing code. This way, the integrity and cohesion of the codebase are preserved, facilitating easier future modifications and extensions.\n\n` + snippet + `mermaid\ngraph LR\nA[Software Development] --\u003e B[Design Patterns]\nB --\u003e C{Cohesion}\nC --\u003e D[Strategy Pattern]\nC --\u003e E[Decorator Pattern]\nD --\u003e F[Encapsulate Changes]\nE --\u003e G[Flexible Extension]\n` + snippet + `\n\n","CreatedAt":"","PublishedAt":"","Reviewed":false,"RelatedPosts":[{"Id":"df81a4ee-dee8-4180-b144-370bf7c90266","Title":"Simplifying Complex Systems with the Facade Pattern"},{"Id":"ec5461e3-6831-4400-8d8f-ba02fee7f337","Title":"Mastering Modularity: The Module Pattern in Depth"},{"Id":"48606d11-f142-4736-948c-d6ddd5374047","Title":"Leveraging Singleton Pattern for Efficient Resource Management"}]}`)),
		blog_models.BlogFromJSON([]byte(`{"Id":"49e6a869-6713-43d1-9431-45f491fea19f","Title":"The Principles of SOLID: Building Blocks for Object-Oriented Design","Body":"The principles of SOLID are foundational concepts in the world of software development, particularly within object-oriented design. SOLID is an acronym that stands for five key principles: Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion. These principles guide developers in creating more maintainable, understandable, and flexible software. The Single Responsibility Principle asserts that a class should only have one reason to change, emphasizing the importance of classes having a single job. The Open/Closed Principle suggests that software entities should be open for extension but closed for modification, allowing systems to grow without altering existing code. Liskov Substitution Principle ensures that objects of a superclass shall be replaceable with objects of subclasses without affecting the correctness of the program. Interface Segregation Principle advocates for many client-specific interfaces instead of one general-purpose interface, and Dependency Inversion Principle highlights that high-level modules should not depend on low-level modules but rather both should depend on abstractions. By adhering to these principles, developers can ensure their code is robust, scalable, and easy to maintain.\n\nUnderstanding and implementing the SOLID principles can significantly improve the quality of object-oriented design and development. Each principle interlocks with the others to provide a comprehensive approach to tackling common software development challenges, such as code rigidity, fragility, and immobility. For instance, applying the Single Responsibility Principle can prevent classes from becoming overly complex and hard to understand, while the Open/Closed Principle encourages the development of systems that are easy to extend without modification, reducing the risk of bugs. Moreover, the Liskov Substitution and Interface Segregation Principles together ensure that the system remains flexible and modular, with well-defined interfaces and interchangeable components. Lastly, Dependency Inversion helps in reducing the coupling between high-level and low-level components of a system, promoting a more decoupled and thus more easily testable and maintainable system. The SOLID principles serve as a roadmap for developers aiming to create high-quality, sustainable, and efficient object-oriented software.\n\nIncorporating the SOLID principles into everyday coding practices is not without its challenges, but the benefits far outweigh the initial learning curve. By fostering a deeper understanding of these principles, developers can build software that is not only functional but also adaptable to changing requirements and technologies. It encourages a thoughtful approach to design, where every aspect of the code is considered in light of these principles, leading to a more cohesive and robust architecture. As the software industry continues to evolve, the principles of SOLID remain timeless guideposts for developing high-quality software that meets the needs of users and businesses alike. Embracing these principles is a step toward mastering object-oriented design and unlocking the full potential of software development.","RelatedPosts":[{"Id":"ca607aef-1d2b-48c7-8b8f-5a0cc3feebc4","Title":"Understanding the Single Responsibility Principle: A Deep Dive"},{"Id":"520de2e8-3192-4f2a-9897-e3d1438c621e","Title":"Exploring the Open/Closed Principle: Extending Software Flexibility"},{"Id":"19ef792c-38e6-499f-9273-7c094177cb89","Title":"The Importance of Interface Segregation in API Design"}]}`)),
		blog_models.BlogFromJSON([]byte(`{"Id":"d0046cbf-4fcc-456e-8e17-1562bb017eba","Title":"behavioural driven development","Body":"### Unlocking the Power of Behavioral Driven Development (BDD) 🚀\n\nBehavioral Driven Development (BDD) is a game-changer in the world of software development, emphasizing collaboration between developers, QA professionals, and non-technical stakeholders. This approach is all about understanding the behavior of an application from the end-user's perspective, making sure that all parties have a clear understanding of the project's goals. By defining behaviors in plain language, everyone involved can ensure that the software truly meets the user's needs. 🎯\n\nAt its core, BDD is about bridging the gap between technical and non-technical team members. Through the use of simple, yet powerful, domain-specific languages like Gherkin, teams can create understandable and executable specifications. These specifications then guide the coding, testing, and maintenance phases, ensuring that every feature is developed with the end user in mind. ✨\n\nImplementing BDD can seem daunting at first, but the benefits are undeniable. Improved communication, clearer understanding of user requirements, and a more focused approach to development are just the tip of the iceberg. When done right, BDD can significantly reduce misunderstandings and rework, leading to faster releases and products that better satisfy user needs. 🛠️\n\nHowever, BDD is more than just a set of practices; it's a mindset shift. Embracing BDD means moving away from a purely task-oriented view of software development to a more collaborative, behavior-focused approach. It encourages teams to ask the right questions, focusing not just on 'how' to implement, but 'why' features should exist in the first place. This leads to more impactful, user-centered software solutions. 🌈\n\nWhether you're a developer, QA professional, or a stakeholder, diving into BDD can profoundly impact how you approach software development. It's about building software that not only works but also delivers real value to users. Ready to take your projects to the next level? It's time to explore Behavioral Driven Development. 🌟","RelatedPosts":[{"Id":"","Title":"Integrating BDD with Agile: A Seamless Approach"},{"Id":"","Title":"The Role of Gherkin in BDD: Simplifying Specifications"},{"Id":"","Title":"Overcoming Common Challenges in Behavioral Driven Development"}]}`)),
		blog_models.BlogFromJSON([]byte(`{"Id":"f08f109a-7a99-4d37-8b3c-bcc2bf169ccd","Title":"doing fizz buzz in javascript step by step","Body":"### Doing Fizz Buzz in JavaScript Step by Step\n\nFizz Buzz is a classic programming challenge often used in interviews or as a coding exercise to check a developer's understanding of basic programming concepts. The task is simple: print numbers from 1 to 100, but for multiples of 3, print 'Fizz', for multiples of 5 print 'Buzz', and for numbers which are multiples of both 3 and 5, print 'FizzBuzz'. Let's dive into how to tackle this in JavaScript! 🚀\n\n**Step 1: The Basic Loop** 🔄\n\nStart by creating a for loop that runs from 1 to 100. This will be the backbone of our program where all logic is applied.\n\n` + snippet + `javascript\nfor (let i = 1; i \u003c= 100; i++) {\n  console.log(i);\n}\n` + snippet + `\n\n**Step 2: Adding Conditions** 🤔\n\nInside the loop, we need to add conditions to check if the current number *i* is a multiple of 3, 5, or both. This is where the fun begins!\n\n` + snippet + `javascript\nfor (let i = 1; i \u003c= 100; i++) {\n  if (i %! (MISSING)=== 0 \u0026\u0026 i %! (MISSING)=== 0) {\n    console.log('FizzBuzz');\n  } else if (i %! (MISSING)=== 0) {\n    console.log('Fizz');\n  } else if (i %! (MISSING)=== 0) {\n    console.log('Buzz');\n  } else {\n    console.log(i);\n  }\n}\n` + snippet + `\n\n**Step 3: Refactoring for Clarity** 🔍\n\nTo make the code cleaner and more readable, you can assign the condition checks to variables or even create a function to handle the logic.\n\n**Step 4: Running Your Code** 🏃\n\nOnce you have your code ready, run it in your JavaScript environment. You should see the numbers from 1 to 100 printed out, with 'Fizz', 'Buzz', and 'FizzBuzz' in the appropriate places.\n\n**Step 5: Enjoy Your Success** 🎉\n\nCongratulations! You've successfully implemented Fizz Buzz in JavaScript. This is a great exercise to understand how to use loops and conditional statements effectively.\n\n**Challenges Ahead** ❓\n\nTry modifying the code to start and end on numbers of your choice, or even implement a GUI to display the results. The possibilities are endless, and it's a great way to practice your JavaScript skills.\n\n**Final Thoughts** 💭\n\nFizz Buzz may seem simple, but it's an excellent way to showcase your understanding of JavaScript basics. Keep experimenting and challenging yourself with different variations of the problem.\n\nHappy Coding! 👨‍💻👩‍💻","RelatedPosts":[{"Id":"","Title":"Understanding JavaScript Loops: For, While, and Beyond"},{"Id":"","Title":"JavaScript Conditional Statements: Making Decisions in Your Code"},{"Id":"","Title":"Creative Coding Challenges to Improve Your JavaScript Skills"}]}`)),
		blog_models.BlogFromJSON([]byte(`{"Id":"d117aef5-1ba3-490b-bdd8-4200f78690b6","Title":"creating a todo list api with nest js","Body":"## Creating a ToDo List API with NestJS 🚀\n\nNestJS, a framework for building efficient, reliable, and scalable server-side applications, has been gaining traction among developers for its out-of-the-box application architecture. If you’re looking to create a ToDo List API, NestJS offers a robust solution that leverages modern JavaScript features and techniques. In this blog post, we’ll walk through the steps to create a basic ToDo List API using NestJS.\n\n### Step 1: Setting Up Your Project 🛠\n\nFirst things first, let’s set up our project environment. You'll need Node.js installed on your computer. Then, install NestJS CLI globally using npm:\n\n` + snippet + `bash\nnpm i -g @nestjs/cli\n` + snippet + `\n\nOnce installed, create a new project:\n\n` + snippet + `bash\nnest new todo-list-api\n` + snippet + `\n\nNavigate into your project directory and you’re all set to start coding!\n\n### Step 2: Create a Task Module 📁\n\nIn NestJS, modules organize related code. For our ToDo List API, we’ll create a Task module. Run the following command:\n\n` + snippet + `bash\nnest generate module tasks\n` + snippet + `\n\nThis command scaffolds a module for us to define task-related logic.\n\n### Step 3: Define a Task Model 📄\n\nCreate a *task.model.ts* file within the tasks module directory. This model will represent the tasks in our database. Define the Task model with properties like id, title, and description:\n\n` + snippet + `typescript\nexport interface Task {\n  id: string;\n  title: string;\n  description: string;\n}\n` + snippet + `\n\n### Step 4: Service and Controller 🚦\n\nNext, generate a service and a controller for the Task module. The service will handle the business logic, while the controller will route requests to the appropriate service methods.\n\n` + snippet + `bash\nnest generate service tasks\nnest generate controller tasks\n` + snippet + `\n\n### Step 5: Implementing CRUD Operations 🛠️\n\nIn the tasks service, implement CRUD operations to manage tasks. Use a simple array to store tasks for now. Here’s how to add a task:\n\n` + snippet + `typescript\nprivate tasks: Task[] = [];\n\naddTask(title: string, description: string): Task {\n  const task: Task = { id: Date.now().toString(), title, description };\n  this.tasks.push(task);\n  return task;\n}\n` + snippet + `\n\n### Step 6: Connecting the Dots 🌐\n\nWith the service and controller in place, wire them up to handle HTTP requests. For adding a task, your controller might look like this:\n\n` + snippet + `typescript\n@PostMapping('/tasks')\nasync addTask(@Body() body): Promise\u003cTask\u003e {\n  return this.tasksService.addTask(body.title, body.description);\n}\n` + snippet + `\n\n### Step 7: Testing Your API 🧪\n\nTest your API using a tool like Postman or Insomnia. Make sure to test each of the CRUD operations you’ve implemented. \n\n### Step 8: Wrapping Up and Running Your Server 🏁\n\nFinally, run your server using the command:\n\n` + snippet + `bash\nnpm run start\n` + snippet + `\n\nVisit *http://localhost:3000/tasks* in your browser or API client to see your ToDo List API in action.\n\n### Conclusion ✨\n\nCreating a ToDo List API with NestJS is straightforward and showcases the power and flexibility of the framework. This guide highlighted the key steps to get a basic API up and running. As you become more familiar with NestJS, you can explore more advanced features, including database integration, authentication, and more.\n\n### Mermaid Diagram: Basic Workflow 📊\n\n` + snippet + `mermaid\ngraph LR\nA[Client] -- Request --\u003e B[Controller]\nB -- Logic --\u003e C[Service]\nC -- Interact --\u003e D[Database]\nD -- Response --\u003e C\nC -- Response --\u003e B\nB -- Response --\u003e A\n` + snippet + `\n\nThis diagram illustrates the basic workflow of a task being processed in our ToDo List API.\n\n","RelatedPosts":[{"Id":"","Title":"Advanced NestJS: Building a Real-time Chat Application"},{"Id":"","Title":"Integrating PostgreSQL with NestJS for Scalable Applications"},{"Id":"","Title":"Securing your NestJS Application with JWT Authentication"}]}`))}

	for _, entity := range entities {
		repo.Create(entity)
	}
}
