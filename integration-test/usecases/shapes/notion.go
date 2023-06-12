package shapes

import (
	"github.com/Permify/permify/pkg/development/file"
)

// NOTION SAMPLE

var InitialNotionShape = file.Shape{
	Schema: `
entity user {}

entity workspace {
    // The owner of the workspace
    relation owner @user
    // Members of the workspace
    relation member @user
    // Guests (users with read-only access) of the workspace
    relation guest @user
    // Bots associated with the workspace
    relation bot @user
    // Admin users who have permission to manage the workspace
    relation admin @user

    // Define permissions for workspace actions
    permission create_page = owner or member or admin
    permission invite_member = owner or admin
    permission view_workspace = owner or member or guest or bot
    permission manage_workspace = owner or admin

    // Define permissions that can be inherited by child entities
    permission read = member or guest or bot or admin
    permission write = owner or admin
}

entity page {
    // The workspace associated with the page
    relation workspace @workspace
     // The user who can write to the page
    relation writer @user
     // The user(s) who can read the page (members of the workspace or guests)
    relation reader @user @workspace#member @workspace#guest

    // Define permissions for page actions
    permission read = reader or workspace.read
    permission write = writer or workspace.write
}

entity database {
    // The workspace associated with the database
    relation workspace @workspace
    // The user who can edit the database
    relation editor @user
    // The user(s) who can view the database (members of the workspace or guests)
    relation viewer @user @workspace#member @workspace#guest

    // Define permissions for database actions
    permission read = viewer or workspace.read
    permission write = editor or workspace.write
    permission create = editor or workspace.write
    permission delete = editor or workspace.write
}

entity block {
    // The page associated with the block
    relation page @page
    // The database associated with the block

    relation database @database
    // The user who can edit the block
    relation editor @user
    // The user(s) who can comment on the block (readers of the parent object)
    relation commenter @user @page#reader

    // Define permissions for block actions
    permission read = database.read or commenter
    permission write = editor or database.write
    permission comment = commenter
}

entity comment {
    // The block associated with the comment
    relation block @block

     // The author of the comment
    relation author @user

    // Define permissions for comment actions
    permission read = block.read
    permission write = author
}

entity template {
   // The workspace associated with the template
    relation workspace @workspace
    // The user who creates the template
    relation creator @user

    // The user(s) who can view the page (members of the workspace or guests)
    relation viewer @user @workspace#member @workspace#guest

    // Define permissions for template actions
    permission read = viewer or workspace.read
    permission write = creator or workspace.write
    permission create = creator or workspace.write
    permission delete = creator or workspace.write
}

entity integration {
    // The workspace associated with the integration
    relation workspace @workspace

    // The owner of the integration
    relation owner @user

    // Define permissions for integration actions
    permission read = workspace.read
    permission write = owner or workspace.write
}
    `,
	Relationships: []string{
		// Assign users to different workspaces:
		"workspace:engineering_team#owner@user:alice",
		"workspace:engineering_team#member@user:bob",
		"workspace:engineering_team#guest@user:charlie",
		"workspace:engineering_team#admin@user:alice",
		"workspace:sales_team#owner@user:david",
		"workspace:sales_team#member@user:eve",
		"workspace:sales_team#guest@user:frank",
		"workspace:sales_team#admin@user:david",

		// Connect pages, databases, and templates to workspaces:
		"page:project_plan#workspace@workspace:engineering_team",
		"page:product_spec#workspace@workspace:engineering_team",
		"database:task_list#workspace@workspace:engineering_team",
		"template:weekly_report#workspace@workspace:sales_team",
		"database:customer_list#workspace@workspace:sales_team",
		"template:marketing_campaign#workspace@workspace:sales_team",

		// Set permissions for pages, databases, and templates:
		"page:project_plan#writer@user:frank",
		"page:project_plan#reader@user:bob",

		"database:task_list#editor@user:alice",
		"database:task_list#viewer@user:bob",

		"template:weekly_report#creator@user:alice",
		"template:weekly_report#viewer@user:bob",

		"page:product_spec#writer@user:david",
		"page:product_spec#reader@user:eve",

		"database:customer_list#editor@user:david",
		"database:customer_list#viewer@user:eve",

		"template:marketing_campaign#creator@user:david",
		"template:marketing_campaign#viewer@user:eve",

		// Set relationships for blocks and comments:
		"block:task_list_1#database@database:task_list",
		"block:task_list_1#editor@user:alice",
		"block:task_list_1#commenter@user:bob",
		"block:task_list_2#database@database:task_list",
		"block:task_list_2#editor@user:alice",
		"block:task_list_2#commenter@user:bob",

		"comment:task_list_1_comment_1#block@block:task_list_1",
		"comment:task_list_1_comment_1#author@user:bob",
		"comment:task_list_1_comment_2#block@block:task_list_1",
		"comment:task_list_1_comment_2#author@user:charlie",
		"comment:task_list_2_comment_1#block@block:task_list_2",
		"comment:task_list_2_comment_1#author@user:bob",
		"comment:task_list_2_comment_2#block@block:task_list_2",
		"comment:task_list_2_comment_2#author@user:charlie",
	},
	Scenarios: []file.Scenario{
		{
			Name:        "Scenario 1",
			Description: "Alice and bob can read to the project plan page",
			Checks: []file.Check{
				{
					ContextualTuples: []string{},
					Entity:           "page:project_plan",
					Subject:          "user:alice",
					Assertions: map[string]bool{
						"read": true,
					},
				},
				{
					ContextualTuples: []string{},
					Entity:           "page:project_plan",
					Subject:          "user:bob",
					Assertions: map[string]bool{
						"read": true,
					},
				},
			},
			EntityFilters: []file.EntityFilter{
				{
					ContextualTuples: []string{
						"page:context#reader@user:bob",
					},
					EntityType: "page",
					Subject:    "user:bob",
					Assertions: map[string][]string{
						"read": {"project_plan", "product_spec", "context"},
					},
				},
			},
			SubjectFilters: []file.SubjectFilter{
				{
					ContextualTuples: []string{},
					Entity:           "page:project_plan",
					SubjectReference: "user",
					Assertions: map[string][]string{
						"read": {"bob", "alice", "charlie"},
					},
				},
			},
		},
		{
			Name:        "Scenario 2",
			Description: "Check if a user who is a guest in a workspace can edit a database",
			Checks: []file.Check{
				{
					ContextualTuples: []string{},
					Entity:           "database:task_list",
					Subject:          "user:frank",
					Assertions: map[string]bool{
						"write": false,
					},
				},
			},
			EntityFilters:  []file.EntityFilter{},
			SubjectFilters: []file.SubjectFilter{},
		},
		{
			Name:        "Scenario 3",
			Description: "Ensure that the owner of a workspace can write to all databases in the workspace",
			Checks:      []file.Check{},
			EntityFilters: []file.EntityFilter{
				{
					ContextualTuples: []string{},
					EntityType:       "database",
					Subject:          "user:alice",
					Assertions: map[string][]string{
						"write": {"task_list"},
					},
				},
			},
			SubjectFilters: []file.SubjectFilter{},
		},
		{
			Name:          "Scenario 4",
			Description:   "Ensure that all members of a workspace can read all pages in the workspace",
			Checks:        []file.Check{},
			EntityFilters: []file.EntityFilter{},
			SubjectFilters: []file.SubjectFilter{
				{
					ContextualTuples: []string{},
					Entity:           "page:project_plan",
					SubjectReference: "user",
					Assertions: map[string][]string{
						"read": {"bob", "alice", "charlie"},
					},
				},
				{
					ContextualTuples: []string{},
					Entity:           "page:product_spec",
					SubjectReference: "user",
					Assertions: map[string][]string{
						"read": {"eve", "bob", "alice", "charlie"},
					},
				},
			},
		},
		{
			Name:        "Scenario 5",
			Description: "Ensure that a user who is not a member of a workspace cannot view the workspace",
			Checks: []file.Check{
				{
					ContextualTuples: []string{},
					Entity:           "workspace:sales_team",
					Subject:          "user:charlie",
					Assertions: map[string]bool{
						"view_workspace": false,
					},
				},
			},
			EntityFilters:  []file.EntityFilter{},
			SubjectFilters: []file.SubjectFilter{},
		},
	},
}
