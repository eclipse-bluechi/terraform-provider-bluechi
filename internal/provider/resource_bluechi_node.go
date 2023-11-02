package provider

import (
	"context"
	"fmt"

	"github.com/eclipse-bluechi/terraform-provider-bluechi/internal/client"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &BlueChiNodeResource{}
var _ resource.ResourceWithImportState = &BlueChiNodeResource{}

func NewBlueChiNodeResource() resource.Resource {
	return &BlueChiNodeResource{}
}

type BlueChiNodeResource struct{}

type BlueChiNodeResourceModel struct {
	Id                types.String            `tfsdk:"id"`
	SSH               BlueChiSSHModel         `tfsdk:"ssh"`
	BlueChiController *BlueChiControllerModel `tfsdk:"bluechi_controller"`
	BlueChiAgent      *BlueChiAgentModel      `tfsdk:"bluechi_agent"`
}

type BlueChiSSHModel struct {
	Host                  types.String `tfsdk:"host"`
	User                  types.String `tfsdk:"user"`
	Password              types.String `tfsdk:"password"`
	PrivateKeyPath        types.String `tfsdk:"private_key_path"`
	AcceptHostKeyInsecure types.Bool   `tfsdk:"accept_host_key_insecure"`
}

type BlueChiControllerModel struct {
	AllowedNodeNames types.Set    `tfsdk:"allowed_node_names"`
	ManagerPort      types.Int64  `tfsdk:"manager_port"`
	LogLevel         types.String `tfsdk:"log_level"`
	LogTarget        types.String `tfsdk:"log_target"`
	LogIsQuiet       types.Bool   `tfsdk:"log_is_quiet"`
	ConfigFile       types.String `tfsdk:"config_file"`
}

func (m BlueChiControllerModel) ToConfig() client.BlueChiControllerConfig {
	cfg := client.BlueChiControllerConfig{}
	m.AllowedNodeNames.ElementsAs(context.Background(), &cfg.AllowedNodeNames, true)
	cfg.ManagerPort = m.ManagerPort.ValueInt64Pointer()
	cfg.LogLevel = m.LogLevel.ValueStringPointer()
	cfg.LogTarget = m.LogTarget.ValueStringPointer()
	cfg.LogIsQuiet = m.LogIsQuiet.ValueBoolPointer()

	return cfg
}

type BlueChiAgentModel struct {
	NodeName          types.String `tfsdk:"node_name"`
	ManagerHost       types.String `tfsdk:"manager_host"`
	ManagerPort       types.Int64  `tfsdk:"manager_port"`
	ManagerAddress    types.String `tfsdk:"manager_address"`
	HeartbeatInterval types.Int64  `tfsdk:"heartbeat_interval"`
	LogLevel          types.String `tfsdk:"log_level"`
	LogTarget         types.String `tfsdk:"log_target"`
	LogIsQuiet        types.Bool   `tfsdk:"log_is_quiet"`
	ConfigFile        types.String `tfsdk:"config_file"`
}

func (m BlueChiAgentModel) ToConfig() client.BlueChiAgentConfig {
	cfg := client.BlueChiAgentConfig{}
	cfg.NodeName = m.NodeName.ValueStringPointer()
	cfg.ManagerHost = m.ManagerHost.ValueStringPointer()
	cfg.ManagerPort = m.ManagerPort.ValueInt64Pointer()
	cfg.ManagerAddress = m.ManagerAddress.ValueStringPointer()
	cfg.HeartbeatInterval = m.HeartbeatInterval.ValueInt64Pointer()
	cfg.LogLevel = m.LogLevel.ValueStringPointer()
	cfg.LogTarget = m.LogTarget.ValueStringPointer()
	cfg.LogIsQuiet = m.LogIsQuiet.ValueBoolPointer()

	return cfg
}

func (r *BlueChiNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (r *BlueChiNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A BlueChi node",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the node",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"ssh": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Connection method to the machine",
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Required:    true,
						Description: "Host of the machine",
						Validators:  []validator.String{},
					},
					"user": schema.StringAttribute{
						Required:    true,
						Description: "User on the machine",
						Validators:  []validator.String{},
					},
					"password": schema.StringAttribute{
						Optional:    true,
						Description: "Password to log in to the machine",
						Validators:  []validator.String{},
					},
					"private_key_path": schema.StringAttribute{
						Optional:    true,
						Description: "Path to the private key used for login",
						Validators:  []validator.String{},
					},
					"accept_host_key_insecure": schema.BoolAttribute{
						Optional:    true,
						Description: "Flag to indicate if host should be validated",
					},
				},
			},
			"bluechi_controller": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "BlueChi controller configuration used on the node",
				Attributes: map[string]schema.Attribute{
					"allowed_node_names": schema.SetAttribute{
						Required:    true,
						ElementType: types.StringType,
						Description: "List of all allowed node names",
						Validators:  []validator.Set{},
					},
					"manager_port": schema.Int64Attribute{
						Optional:    true,
						Description: "Port the manager is listening on",
						Validators:  []validator.Int64{},
					},
					"log_level": schema.StringAttribute{
						Optional:    true,
						Description: "Log level used by BlueChi controller",
						Validators:  []validator.String{},
					},
					"log_target": schema.StringAttribute{
						Optional:    true,
						Description: "Log target used by BlueChi controller",
						Validators:  []validator.String{},
					},
					"log_is_quiet": schema.BoolAttribute{
						Optional:    true,
						Description: "Flag to indicate if logs are written",
					},
					"config_file": schema.StringAttribute{
						Computed:    true,
						Description: "The bluechi controller configuration file on the system",
					},
				},
			},
			"bluechi_agent": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "BlueChi agent configuration used on the node",
				Attributes: map[string]schema.Attribute{
					"node_name": schema.StringAttribute{
						Required:    true,
						Description: "Name of the BlueChi agent",
						Validators:  []validator.String{},
					},
					"manager_host": schema.StringAttribute{
						Required:    true,
						Description: "Host of the manager to connect to",
						Validators:  []validator.String{},
					},
					"manager_port": schema.Int64Attribute{
						Required:    true,
						Description: "Port of the manager to connect to",
						Validators:  []validator.Int64{},
					},
					"manager_address": schema.StringAttribute{
						Optional:    true,
						Description: "Address of the manager to connect to. Replaces host and port.",
						Validators:  []validator.String{},
					},
					"heartbeat_interval": schema.Int64Attribute{
						Optional:    true,
						Description: "The interval in ms in which the connection is tested",
						Validators:  []validator.Int64{},
					},
					"log_level": schema.StringAttribute{
						Optional:    true,
						Description: "Log level used by BlueChi controller",
						Validators:  []validator.String{},
					},
					"log_target": schema.StringAttribute{
						Optional:    true,
						Description: "Log target used by BlueChi controller",
						Validators:  []validator.String{},
					},
					"log_is_quiet": schema.BoolAttribute{
						Optional:    true,
						Description: "Flag to indicate if logs are written",
					},
					"config_file": schema.StringAttribute{
						Computed:    true,
						Description: "The bluechi agent configuration file on the system",
					},
				},
			},
		},
	}
}

func (r *BlueChiNodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (r *BlueChiNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BlueChiNodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to read model")
		return
	}

	sshClient, errDiag := setupSSHClient(data.SSH)
	if errDiag != nil {
		tflog.Error(ctx, "Failed to connect via SSH")
		resp.Diagnostics.AddError(errDiag.Summary(), errDiag.Detail())
		return
	}
	defer sshClient.Disconnect()

	id, err := uuid.GenerateUUID()
	if err != nil {
		tflog.Error(ctx, "Failed to generate UUID for node")
		resp.Diagnostics.AddError("Failed to generate UUID for node", err.Error())
		return
	}
	data.Id = types.StringValue(id)

	ctrlConf := data.BlueChiController
	agentConf := data.BlueChiAgent

	err = sshClient.InstallBlueChi(ctrlConf != nil, agentConf != nil)
	if err != nil {
		tflog.Error(ctx, "Failed to install BlueChi")
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to install BlueChi: %v", err), err.Error())
		return
	}

	if ctrlConf != nil {
		ctrlConfFile := assembleConfigFileName("ctrl")
		err := sshClient.CreateControllerConfig(ctrlConfFile, data.BlueChiController.ToConfig())
		if err != nil {
			tflog.Error(ctx, "Failed to create controller config")
			resp.Diagnostics.AddError("Failed to create controller config", err.Error())
			return
		}
		data.BlueChiController.ConfigFile = types.StringValue(ctrlConfFile)

		err = sshClient.RestartBlueChiController()
		if err != nil {
			tflog.Error(ctx, "Failed to start controller service")
			resp.Diagnostics.AddError("Failed to start controller service", err.Error())
			return
		}
	}

	if agentConf != nil {
		agentConfFile := assembleConfigFileName("agent")
		err := sshClient.CreateAgentConfig(agentConfFile, data.BlueChiAgent.ToConfig())
		if err != nil {
			tflog.Error(ctx, "Failed to create agent config")
			resp.Diagnostics.AddError("Failed to create agent config", err.Error())
			return
		}
		data.BlueChiAgent.ConfigFile = types.StringValue(agentConfFile)

		err = sshClient.RestartBlueChiAgent()
		if err != nil {
			tflog.Error(ctx, "Failed to start agent service")
			resp.Diagnostics.AddError("Failed to start agent service", err.Error())
			return
		}
	}

	tflog.Trace(ctx, "Setup BlueChi on machine completed")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlueChiNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BlueChiNodeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: implement read

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlueChiNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BlueChiNodeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sshClient, errDiag := setupSSHClient(data.SSH)
	if errDiag != nil {
		tflog.Error(ctx, "Failed to create and connect via SSH")
		resp.Diagnostics.AddError(errDiag.Summary(), errDiag.Detail())
		return
	}
	defer sshClient.Disconnect()

	ctrlConf := data.BlueChiController
	if ctrlConf != nil {
		err := sshClient.CreateControllerConfig(
			data.BlueChiController.ConfigFile.ValueString(),
			data.BlueChiController.ToConfig(),
		)
		if err != nil {
			tflog.Error(ctx, "Failed to update controller config")
			resp.Diagnostics.AddError("Failed to update controller config", err.Error())
			return
		}

		err = sshClient.RestartBlueChiController()
		if err != nil {
			tflog.Error(ctx, "Failed to start controller service")
			resp.Diagnostics.AddError("Failed to start controller service", err.Error())
			return
		}
	}

	agentConf := data.BlueChiAgent
	if agentConf != nil {
		err := sshClient.CreateAgentConfig(
			data.BlueChiAgent.ConfigFile.ValueString(),
			data.BlueChiAgent.ToConfig(),
		)
		if err != nil {
			tflog.Error(ctx, "Failed to update agent config")
			resp.Diagnostics.AddError("Failed to update agent config", err.Error())
			return
		}

		err = sshClient.RestartBlueChiAgent()
		if err != nil {
			tflog.Error(ctx, "Failed to start agent service")
			resp.Diagnostics.AddError("Failed to start agent service", err.Error())
			return
		}
	}

	tflog.Trace(ctx, "Setup BlueChi on machine updated")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BlueChiNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BlueChiNodeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshClient, errDiag := setupSSHClient(data.SSH)
	if sshClient == nil {
		tflog.Error(ctx, "Failed to connect via SSH")
		resp.Diagnostics.AddError(errDiag.Summary(), errDiag.Detail())
		return
	}
	defer sshClient.Disconnect()

	ctrlConf := data.BlueChiController
	if ctrlConf != nil {
		err := sshClient.RemoveControllerConfig(ctrlConf.ConfigFile.ValueString())
		if err != nil {
			tflog.Error(ctx, "Failed to remove controller config")
			resp.Diagnostics.AddError("Failed to remove controller config", err.Error())
			return
		}

		err = sshClient.StopBlueChiController()
		if err != nil {
			tflog.Error(ctx, "Failed to stop controller service")
			resp.Diagnostics.AddError("Failed to stop controller service", err.Error())
			return
		}
	}

	agentConf := data.BlueChiAgent
	if agentConf != nil {
		err := sshClient.RemoveAgentConfig(agentConf.ConfigFile.ValueString())
		if err != nil {
			tflog.Error(ctx, "Failed to remove agent config")
			resp.Diagnostics.AddError("Failed to remove agent config", err.Error())
			return
		}

		err = sshClient.StopBlueChiAgent()
		if err != nil {
			tflog.Error(ctx, "Failed to stop agent service")
			resp.Diagnostics.AddError("Failed to stop agent service", err.Error())
			return
		}
	}
}

func (r *BlueChiNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setupSSHClient(sshModel BlueChiSSHModel) (client.Client, *diag.ErrorDiagnostic) {
	sshClient := client.NewSSHClient(
		sshModel.Host.ValueString(),
		sshModel.User.ValueString(),
		sshModel.Password.ValueString(),
		sshModel.PrivateKeyPath.ValueString(),
		sshModel.AcceptHostKeyInsecure.ValueBool(),
	)

	if err := sshClient.Connect(); err != nil {
		errSummary := fmt.Sprintf("Failed to connect to '%s'", sshModel.Host.ValueString())
		diagnostic := diag.NewErrorDiagnostic(errSummary, err.Error())
		return nil, &diagnostic
	}

	return sshClient, nil
}

func assembleConfigFileName(suffix string) string {
	return fmt.Sprintf("ZZZ-%s.conf", suffix)
}
