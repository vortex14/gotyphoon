package graph
//
//import (
//	"fmt"
//	"github.com/vortex14/gotyphoon/elements/forms"
//	"github.com/vortex14/gotyphoon/elements/models/label"
//	ghvzExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
//	"github.com/vortex14/gotyphoon/interfaces"
//)
//
//type TyphoonServer struct {
//	*forms.TyphoonServer
//
//	BuildGraph      bool
//	Graph           interfaces.GraphInterface
//
//}
//
//
//func (s *TyphoonServer) InitGraph() interfaces.ServerInterface {
//	s.Graph = (&ghvzExt.Graph{
//		MetaInfo: &label.MetaInfo{
//			Name: fmt.Sprintf("Graph of %s",s.Name),
//		},
//		Layout: ghvzExt.LAYOUTCirco,
//	}).Init()
//	return s
//}
//
//func (s *TyphoonServer) GetGraph() interfaces.GraphInterface {
//	return s.Graph
//}
//
//func (s *TyphoonServer) AddNewGraphResource(newResource interfaces.ResourceGraphInterface)  {
//	if s.Graph != nil {
//		subGraph := s.Graph.AddSubGraph(&interfaces.GraphOptions{
//			Name:      newResource.GetName(),
//			Label:     newResource.GetName(),
//			IsCluster: true,
//		})
//		s.LOG.Debug(fmt.Sprintf("init subGraph for %s", newResource.GetName()), subGraph)
//		newResource.SetGraph(subGraph)
//	} else {
//		s.LOG.Error("not found server graph. ",newResource.GetPath(), newResource.GetName())
//	}
//}