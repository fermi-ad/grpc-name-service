package nameserver;
import com.google.protobuf.Empty;

import io.quarkus.hibernate.orm.panache.PanacheQuery;
import io.quarkus.panache.common.Parameters;
import jakarta.transaction.Transactional;
import nameserver.model.Node;
import java.util.List;
import java.util.Arrays;
import java.util.ArrayList;

/*Node API implementation */
public class NodeController {
    
    public static Node getNodeByHostname(String hostname) {        
        Node node = Node.find("hostname", hostname).firstResult();
        if (node == null) {
            throw new RuntimeException("Node with hostname " + hostname + " not found");
        }
        return node;
    }
    
    static proto.Node.Builder createNodeBuilder(Node node) {
        var nodeBuilder = proto.Node.newBuilder();
        
        nodeBuilder.setHostname(node.getHostname())
                .setDescription(node.getDescription())
                .setIpAddress(node.getIpAddress())
                .setLocationName(node.getLocation().getName());
        return nodeBuilder;
    }

    @Transactional
    public static proto.NodeResponse createNode(proto.CreateNodeRequest request) {        
        var reqNode = request.getNode();
        var location = LocationController.getLocationByName(reqNode.getLocationName());
        Node node = new Node(reqNode.getHostname(),
                             reqNode.getIpAddress(),
                             location,
                             reqNode.getDescription());
        node.persist();
        var nodeBuilder = createNodeBuilder(node);        
        return proto.NodeResponse.newBuilder().setNode(nodeBuilder).build();
    }

    @Transactional
    public static proto.NodeResponse updateNode(proto.UpdateNodeRequest request) {
        var reqNode = request.getNode();
        Node node = getNodeByHostname(reqNode.getHostname());        
        if (reqNode.getDescription() != "") {
            node.setDescription(reqNode.getDescription());
        }
        if (reqNode.getIpAddress() != "") {
            node.setIpAddress(reqNode.getIpAddress());
        }
        if (reqNode.getLocationName() != "") {
            node.setLocation(LocationController.getLocationByName(reqNode.getLocationName()));
        }
        node.persist();

        var nodeBuilder = createNodeBuilder(node);        
        return proto.NodeResponse.newBuilder().setNode(nodeBuilder).build();
    }

    public static proto.NodeListResponse listNodes(proto.ListNodesRequest request) {
        var queryParams = new ArrayList<List<String>>();
        if (!request.getHostname().isEmpty()) {
            queryParams.add(Arrays.asList("hostname",request.getHostname()));
        }
        if (!request.getIpAddress().isEmpty()) {
            queryParams.add(Arrays.asList("ipAddress",request.getIpAddress()));
        }
        if (!request.getLocationName().isEmpty()) {
            queryParams.add(Arrays.asList("location.name",request.getLocationName()));
        }
        Parameters parameters = new Parameters();   
            
        var queryString = ControllerUtil.createFindParam(queryParams, parameters);        
        PanacheQuery<Node> query;
        query = Node.find(queryString, parameters);
        
        var pag = ControllerUtil.readPagination(request.getPagination());                
        List<Node> nodes = query.page(pag[0] - 1, pag[1]).list();        
        proto.NodeListResponse.Builder response = proto.NodeListResponse.newBuilder();
        for (Node node : nodes) {
            var nodeBuilder = createNodeBuilder(node);            
            response.addNodes(nodeBuilder.build());
        }
        response.setPagination(ControllerUtil.createPaginationResponse(pag, query.count()));
        return response.build();
    }

    @Transactional
    public static Empty deleteNode(proto.DeleteNodeRequest request) {        
        Node.delete("hostname", request.getHostname());
        return Empty.newBuilder().build();
    }     

    public static proto.NodeResponse getNode(proto.GetNodeRequest request) {
        var node = getNodeByHostname(request.getHostname());
        var nodeBuilder = createNodeBuilder(node);
        return proto.NodeResponse.newBuilder().setNode(nodeBuilder).build();
    }
 }
