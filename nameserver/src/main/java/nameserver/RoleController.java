package nameserver;
import com.google.protobuf.Empty;

import jakarta.transaction.Transactional;
import nameserver.model.Role;
import java.util.List;

/*Role API implementation */
public class RoleController {
    public static Role getRoleByName(String name) {        
        Role role = Role.find("name", name).firstResult();
        if (role == null) {
            throw new RuntimeException("Role with name " + name + " not found");
        }
        return role;
    }

    static proto.Role.Builder createRoleBuilder(Role role) {
        var roleBuilder = proto.Role.newBuilder();
        
        roleBuilder.setName(role.getName())
                .setDescription(role.getDescription());
        return roleBuilder;
    }

    @Transactional
    public static proto.RoleResponse createRole(proto.CreateRoleRequest request) {        
        var reqLocType = request.getRole();        
        var role = new Role(reqLocType.getName(), reqLocType.getDescription());
        try {
            role.persist();
        } catch (Exception e) {
            throw new RuntimeException("Error creating Role: " + e.getMessage());
        }
        var resLocTypeBuilder = createRoleBuilder(role);
        return proto.RoleResponse.newBuilder().setRole(resLocTypeBuilder).build();
    }

    @Transactional
    public static proto.RoleResponse updateRole(proto.UpdateRoleRequest request) {
        var reqLocType = request.getRole();        
        var name = reqLocType.getName();
        var role = getRoleByName(name);
        if (reqLocType.getDescription() != "") {
            role.setDescription(reqLocType.getDescription());
        }
        try {
            role.persist();
        } catch (Exception e) {
            throw new RuntimeException("Error updating Role: " + e.getMessage());
        }
        
        var resLocTypeBuilder = createRoleBuilder(role);
        return proto.RoleResponse.newBuilder().setRole(resLocTypeBuilder).build();        
    }

    public static proto.RoleListResponse listRoles(Empty request) {
        var responseBuilder = proto.RoleListResponse.newBuilder();
        List<Role> roles = Role.listAll();
        for (Role l : roles) {
            var lbuilder = createRoleBuilder(l);
            responseBuilder.addRoles(lbuilder.build());
        }
        return responseBuilder.build();
    }

    @Transactional
    public static Empty deleteRole(proto.DeleteRoleRequest request) {        
        Role.delete("name", request.getName());
        return Empty.newBuilder().build();
    }
}
