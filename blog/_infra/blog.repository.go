package blog_infra

import (
	blog_models "ioprodz/blog/_models"
	"ioprodz/common/policies"
)

type BlogMemoryRepository struct {
	list []blog_models.Blog
}

func (repo *BlogMemoryRepository) Create(qna blog_models.Blog) {
	repo.list = append(repo.list, qna)
}

func (repo *BlogMemoryRepository) List() []blog_models.Blog {
	return repo.list
}

func (repo *BlogMemoryRepository) Get(id string) (blog_models.Blog, error) {
	for _, obj := range repo.list {
		if obj.Id == id {
			return obj, nil
		}
	}
	return blog_models.Blog{}, &policies.StorageError{Message: "Element not found by id: " + id}
}

func CreateBlogRepo() *BlogMemoryRepository {
	repo := &BlogMemoryRepository{list: make([]blog_models.Blog, 0)}
	repo.seed()
	return repo
}

func (repo *BlogMemoryRepository) seed() {
	repo.list = []blog_models.Blog{blog_models.BlogFromJSON([]byte(`{"Id":"523cfd06-6283-4d3a-937a-5dc68ffcb9ed","Title":"Understanding coupling and cohesion","Body":"In the realm of software engineering, two fundamental concepts that stand as pillars for designing robust, maintainable, and scalable systems are coupling and cohesion. Coupling refers to the degree of direct knowledge or reliance one class or module has on another. This means that changes in one module could directly affect another module if the coupling is tight, leading to a fragile system architecture where changes can cascade and cause unforeseen issues. Conversely, low coupling or loose coupling aims for modules that are more independent, making the system more modular and easier to modify or extend. \n\n Cohesion, on the other hand, deals with how closely related and focused the responsibilities of a single module are. High cohesion within a module means that its functions are closely related to each other and perform a well-defined task, making the system more understandable and easier to maintain. It encourages the breaking down of a system into manageable pieces, where each piece is responsible for a single, focused task or group of related tasks. High cohesion and low coupling are often seen as ideals in software design, leading to systems that are easier to debug, test, and maintain over time. \n\n Understanding and applying the principles of coupling and cohesion can significantly impact the quality of software systems. By striving for low coupling, developers ensure that modules can be developed, tested, and debugged in isolation, reducing complexity and dependency issues. Similarly, by aiming for high cohesion within modules, developers can create more intuitive and manageable codebases. These principles guide the design of software that is not only functional but also adaptable and resilient to change, qualities that are essential in the fast-paced and ever-evolving field of software development.","RelatedPosts":[{"Id":"49e6a869-6713-43d1-9431-45f491fea19f","Title":"The Principles of SOLID: Building Blocks for Object-Oriented Design"},{"Id":"82a40beb-bb7e-4874-8bd3-ffac983e555c","Title":"Refactoring Techniques for Reducing Coupling in Your Code"},{"Id":"767bdd8b-9972-4912-beb1-948cfcad4ae2","Title":"Design Patterns: Strategies for Enhancing Cohesion"}]}`)),
		blog_models.BlogFromJSON([]byte(`{"Id":"49e6a869-6713-43d1-9431-45f491fea19f","Title":"The Principles of SOLID: Building Blocks for Object-Oriented Design","Body":"The principles of SOLID are foundational concepts in the world of software development, particularly within object-oriented design. SOLID is an acronym that stands for five key principles: Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion. These principles guide developers in creating more maintainable, understandable, and flexible software. The Single Responsibility Principle asserts that a class should only have one reason to change, emphasizing the importance of classes having a single job. The Open/Closed Principle suggests that software entities should be open for extension but closed for modification, allowing systems to grow without altering existing code. Liskov Substitution Principle ensures that objects of a superclass shall be replaceable with objects of subclasses without affecting the correctness of the program. Interface Segregation Principle advocates for many client-specific interfaces instead of one general-purpose interface, and Dependency Inversion Principle highlights that high-level modules should not depend on low-level modules but rather both should depend on abstractions. By adhering to these principles, developers can ensure their code is robust, scalable, and easy to maintain.\n\nUnderstanding and implementing the SOLID principles can significantly improve the quality of object-oriented design and development. Each principle interlocks with the others to provide a comprehensive approach to tackling common software development challenges, such as code rigidity, fragility, and immobility. For instance, applying the Single Responsibility Principle can prevent classes from becoming overly complex and hard to understand, while the Open/Closed Principle encourages the development of systems that are easy to extend without modification, reducing the risk of bugs. Moreover, the Liskov Substitution and Interface Segregation Principles together ensure that the system remains flexible and modular, with well-defined interfaces and interchangeable components. Lastly, Dependency Inversion helps in reducing the coupling between high-level and low-level components of a system, promoting a more decoupled and thus more easily testable and maintainable system. The SOLID principles serve as a roadmap for developers aiming to create high-quality, sustainable, and efficient object-oriented software.\n\nIncorporating the SOLID principles into everyday coding practices is not without its challenges, but the benefits far outweigh the initial learning curve. By fostering a deeper understanding of these principles, developers can build software that is not only functional but also adaptable to changing requirements and technologies. It encourages a thoughtful approach to design, where every aspect of the code is considered in light of these principles, leading to a more cohesive and robust architecture. As the software industry continues to evolve, the principles of SOLID remain timeless guideposts for developing high-quality software that meets the needs of users and businesses alike. Embracing these principles is a step toward mastering object-oriented design and unlocking the full potential of software development.","RelatedPosts":[{"Id":"ca607aef-1d2b-48c7-8b8f-5a0cc3feebc4","Title":"Understanding the Single Responsibility Principle: A Deep Dive"},{"Id":"520de2e8-3192-4f2a-9897-e3d1438c621e","Title":"Exploring the Open/Closed Principle: Extending Software Flexibility"},{"Id":"19ef792c-38e6-499f-9273-7c094177cb89","Title":"The Importance of Interface Segregation in API Design"}]}`))}
}
